//go:generate go run -mod=mod go.uber.org/mock/mockgen -typed -source=$GOFILE -destination=../repository/mock/client.go -package=repository -build_flags=-mod=mod
package service

import (
	"context"

	"neoway/pkg/entity"
)

// SearchLimit is the limit of clients to search at once.
const SearchLimit = 1000

// ClientRepository is a repository that handles clients.
type ClientRepository interface {
	CreateMulti(ctx context.Context, clients []*entity.Client) error
	GetPaginated(ctx context.Context, limit, offset int) ([]*entity.Client, error)
	UpdateMultiDocumentsValidity(ctx context.Context, clients []*entity.Client) error
}

type client struct {
	repo ClientRepository
}

// NewClient creates a new client service.
func NewClient(repo ClientRepository) *client {
	return &client{
		repo: repo,
	}
}

// CreateMulti creates multiple clients in the database.
func (c *client) CreateMulti(ctx context.Context, clients []*entity.Client) error {
	return c.repo.CreateMulti(ctx, clients)
}

// NormalizeClients gets all clients from the database and normalizes them.
// It processes the clients in batches of SearchLimit.
func (c *client) NormalizeClients(ctx context.Context) error {
	for offset := 0; ; offset += SearchLimit {
		clients, err := c.repo.GetPaginated(ctx, SearchLimit, offset)
		if err != nil {
			return err
		}

		if len(clients) == 0 {
			break
		}

		for _, client := range clients {
			client.ValidateDocuments()
		}

		if err := c.repo.UpdateMultiDocumentsValidity(ctx, clients); err != nil {
			return err
		}
	}

	return nil
}
