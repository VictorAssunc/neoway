package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"neoway/pkg/entity"
	"neoway/pkg/lib/pointer"
)

func TestNewClient(t *testing.T) {
	require.NotNil(t, NewClient(nil))
}

func Test_client_CreateMulti(t *testing.T) {
	ctx := context.Background()

	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := NewClient(db)

	clients := []*entity.Client{
		{
			CPF:               "12345678900",
			Private:           true,
			Incomplete:        true,
			LastOrderDate:     pointer.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			AverageTicket:     pointer.New(100.00),
			LastOrderTicket:   pointer.New(200.00),
			MostFrequentStore: pointer.New("12345678900001"),
			LastOrderStore:    pointer.New("12345678900001"),
		},
		{
			CPF:               "98765432100",
			Private:           false,
			Incomplete:        false,
			LastOrderDate:     nil,
			AverageTicket:     nil,
			LastOrderTicket:   nil,
			MostFrequentStore: nil,
			LastOrderStore:    nil,
		},
	}

	placeholders := "($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11),($12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)"
	query := regexp.QuoteMeta(fmt.Sprintf(repo.queries["insert-multi"], placeholders))

	t.Run("success", func(t *testing.T) {
		dbmock.
			ExpectPrepare(query).
			ExpectExec().
			WithArgs(
				"12345678900", true, true, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100.00, 200.00, "12345678900001", "12345678900001", false, false, false,
				"98765432100", false, false, nil, nil, nil, nil, nil, false, false, false,
			).
			WillReturnResult(sqlmock.NewResult(1, 2))

		require.NoError(t, repo.CreateMulti(ctx, clients))
	})

	t.Run("error", func(t *testing.T) {
		t.Run("prepare query", func(t *testing.T) {
			expectedErr := errors.New("error")
			dbmock.ExpectPrepare(query).WillReturnError(expectedErr)

			require.ErrorIs(t, repo.CreateMulti(ctx, clients), expectedErr)
		})

		t.Run("exec query", func(t *testing.T) {
			expectedErr := errors.New("error")
			dbmock.
				ExpectPrepare(query).
				ExpectExec().
				WillReturnError(expectedErr)

			require.ErrorIs(t, repo.CreateMulti(ctx, clients), expectedErr)

		})
	})
}

func Test_client_GetPaginated(t *testing.T) {
	ctx := context.Background()

	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := NewClient(db)

	clients := []*entity.Client{
		{
			ID:                1,
			CPF:               "12345678900",
			Private:           true,
			Incomplete:        true,
			LastOrderDate:     pointer.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			AverageTicket:     pointer.New(100.00),
			LastOrderTicket:   pointer.New(200.00),
			MostFrequentStore: pointer.New("12345678900001"),
			LastOrderStore:    pointer.New("12345678900001"),
		},
		{
			ID:                2,
			CPF:               "98765432100",
			Private:           false,
			Incomplete:        false,
			LastOrderDate:     nil,
			AverageTicket:     nil,
			LastOrderTicket:   nil,
			MostFrequentStore: nil,
			LastOrderStore:    nil,
		},
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "cpf", "private", "incomplete", "last_order_date", "average_ticket", "last_order_ticket", "most_frequent_store", "last_order_store"}).
			AddRow(1, "12345678900", true, true, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 100.00, 200.00, "12345678900001", "12345678900001").
			AddRow(2, "98765432100", false, false, nil, nil, nil, nil, nil)
		dbmock.ExpectQuery(regexp.QuoteMeta(repo.queries["get-paginated"])).WithArgs(2, 0).WillReturnRows(rows)

		result, err := repo.GetPaginated(ctx, 2, 0)
		require.NoError(t, err)
		require.Equal(t, clients, result)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("error")
		dbmock.ExpectQuery(regexp.QuoteMeta(repo.queries["get-paginated"])).WithArgs(2, 0).WillReturnError(expectedErr)

		result, err := repo.GetPaginated(ctx, 2, 0)
		require.Error(t, err, expectedErr)
		require.Nil(t, result)
	})
}

func Test_client_UpdateMultiDocumentsValidity(t *testing.T) {
	ctx := context.Background()

	db, dbmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := NewClient(db)

	clients := []*entity.Client{
		{
			ID:                     1,
			CPF:                    "12345678900",
			ValidCPF:               true,
			ValidMostFrequentStore: true,
			ValidLastOrderStore:    true,
		},
		{
			ID:                     2,
			CPF:                    "98765432100",
			ValidCPF:               false,
			ValidMostFrequentStore: false,
			ValidLastOrderStore:    false,
		},
	}

	placeholders := "($1, $2, $3, $4, $5, CURRENT_TIMESTAMP),($6, $7, $8, $9, $10, CURRENT_TIMESTAMP)"
	query := regexp.QuoteMeta(fmt.Sprintf(repo.queries["update-multi-documents-validity"], placeholders))

	t.Run("success", func(t *testing.T) {
		dbmock.
			ExpectPrepare(query).
			ExpectExec().
			WithArgs(
				1, "12345678900", true, true, true,
				2, "98765432100", false, false, false,
			).
			WillReturnResult(sqlmock.NewResult(1, 2))

		require.NoError(t, repo.UpdateMultiDocumentsValidity(ctx, clients))
	})

	t.Run("error", func(t *testing.T) {
		t.Run("prepare query", func(t *testing.T) {
			expectedErr := errors.New("error")
			dbmock.ExpectPrepare(query).WillReturnError(expectedErr)

			require.ErrorIs(t, repo.UpdateMultiDocumentsValidity(ctx, clients), expectedErr)
		})

		t.Run("exec query", func(t *testing.T) {
			expectedErr := errors.New("error")
			dbmock.
				ExpectPrepare(query).
				ExpectExec().
				WillReturnError(expectedErr)

			require.ErrorIs(t, repo.UpdateMultiDocumentsValidity(ctx, clients), expectedErr)

		})
	})
}
