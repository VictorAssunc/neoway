//go:generate go run -mod=mod go.uber.org/mock/mockgen -typed -source=$GOFILE -destination=../service/mock/client.go -package=service -build_flags=-mod=mod
package handler

import (
	"context"

	"neoway/pkg/entity"
	"neoway/pkg/lib/regex"
)

// DefaultBulkSize is the default bulk size to create clients.
const DefaultBulkSize = 1000

// ClientService is a service that handles clients.
type ClientService interface {
	CreateMulti(ctx context.Context, clients []*entity.Client) error
	NormalizeClients(ctx context.Context) error
}

type client struct {
	clientService ClientService
}

// NewClient creates a new client handler.
func NewClient(clientService ClientService) *client {
	return &client{
		clientService: clientService,
	}
}

// CreateClientsFromLines creates clients from lines and stores them in the database.
func (h *client) CreateClientsFromLines(ctx context.Context, lines []string) error {
	bulkSize := DefaultBulkSize
	if len(lines) < DefaultBulkSize {
		bulkSize = len(lines)
	}

	clients := make([]*entity.Client, 0, bulkSize)
	for i, line := range lines {
		client, err := entity.NewClient(regex.SpaceSequence.Split(line, -1))
		if err != nil {
			return err
		}

		clients = append(clients, client)
		if len(clients) >= bulkSize || i == len(lines)-1 {
			if err := h.clientService.CreateMulti(ctx, clients); err != nil {
				return err
			}

			clients = clients[:0]
		}
	}

	return nil
}

// NormalizeClients normalizes clients.
func (h *client) NormalizeClients(ctx context.Context) error {
	return h.clientService.NormalizeClients(ctx)
}
