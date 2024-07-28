package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"strings"

	"github.com/nleof/goyesql"

	"neoway/pkg/entity"
)

var (
	//go:embed query/client.sql
	queriesFile       []byte
	insertFieldsCount = 11
	updateFieldsCount = 5
)

type client struct {
	db      *sql.DB
	queries goyesql.Queries
}

// NewClient creates a new client repository.
func NewClient(db *sql.DB) *client {
	return &client{
		db:      db,
		queries: goyesql.MustParseBytes(queriesFile),
	}
}

// CreateMulti creates multiple clients in the database.
func (c *client) CreateMulti(ctx context.Context, clients []*entity.Client) error {
	placeholders := make([]string, 0, len(clients))
	values := make([]interface{}, 0, len(clients)*insertFieldsCount)
	for i, client := range clients {
		client := client
		placeholders = append(placeholders,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				i*insertFieldsCount+1,
				i*insertFieldsCount+2,
				i*insertFieldsCount+3,
				i*insertFieldsCount+4,
				i*insertFieldsCount+5,
				i*insertFieldsCount+6,
				i*insertFieldsCount+7,
				i*insertFieldsCount+8,
				i*insertFieldsCount+9,
				i*insertFieldsCount+10,
				i*insertFieldsCount+11,
			),
		)
		values = append(values,
			client.CPF,
			client.Private,
			client.Incomplete,
			client.LastOrderDate,
			client.AverageTicket,
			client.LastOrderTicket,
			client.MostFrequentStore,
			client.LastOrderStore,
			client.ValidCPF,
			client.ValidMostFrequentStore,
			client.ValidLastOrderStore,
		)
	}

	query := fmt.Sprintf(c.queries["insert-multi"], strings.Join(placeholders, ","))
	stmt, err := c.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	return err
}

// GetPaginated returns a paginated list of clients from the database.
func (c *client) GetPaginated(ctx context.Context, limit, offset int) ([]*entity.Client, error) {
	rows, err := c.db.QueryContext(ctx, c.queries["get-paginated"], limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*entity.Client
	for rows.Next() {
		var client entity.Client
		if err := rows.Scan(
			&client.ID,
			&client.CPF,
			&client.Private,
			&client.Incomplete,
			&client.LastOrderDate,
			&client.AverageTicket,
			&client.LastOrderTicket,
			&client.MostFrequentStore,
			&client.LastOrderStore,
		); err != nil {
			return nil, err
		}

		clients = append(clients, &client)
	}

	return clients, nil
}

// UpdateMultiDocumentsValidity updates the validity of multiple clients documents in the database.
func (c *client) UpdateMultiDocumentsValidity(ctx context.Context, clients []*entity.Client) error {
	placeholders := make([]string, 0, len(clients))
	values := make([]interface{}, 0, len(clients)*updateFieldsCount)
	for i, client := range clients {
		client := client
		placeholders = append(placeholders,
			fmt.Sprintf(
				"($%d, $%d, $%d, $%d, $%d, CURRENT_TIMESTAMP)",
				i*updateFieldsCount+1,
				i*updateFieldsCount+2,
				i*updateFieldsCount+3,
				i*updateFieldsCount+4,
				i*updateFieldsCount+5,
			),
		)
		values = append(values,
			client.ID,
			client.CPF,
			client.ValidCPF,
			client.ValidMostFrequentStore,
			client.ValidLastOrderStore,
		)
	}

	query := fmt.Sprintf(c.queries["update-multi-documents-validity"], strings.Join(placeholders, ","))
	stmt, err := c.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	return err
}
