package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"neoway/pkg/entity"
	"neoway/pkg/lib/pointer"
	service "neoway/pkg/service/mock"
)

func TestNewClient(t *testing.T) {
	require.NotNil(t, NewClient(nil))
}

func Test_client_CreateClientsFromLines(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientService := service.NewMockClientService(ctrl)
	handler := NewClient(clientService)

	lines := []string{"12345678900 1 1 2020-01-01 100.00 200.00 12345678900001 12345678900001", "98765432100 1 1 2020-01-01 100.00 200.00 12345678900001 12345678900001"}
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
			Private:           true,
			Incomplete:        true,
			LastOrderDate:     pointer.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
			AverageTicket:     pointer.New(100.00),
			LastOrderTicket:   pointer.New(200.00),
			MostFrequentStore: pointer.New("12345678900001"),
			LastOrderStore:    pointer.New("12345678900001"),
		},
	}

	t.Run("success", func(t *testing.T) {
		clientService.EXPECT().CreateMulti(ctx, clients).Return(nil)

		require.NoError(t, handler.CreateClientsFromLines(ctx, lines))
	})

	t.Run("error", func(t *testing.T) {
		t.Run("new client", func(t *testing.T) {
			line := []string{"12345678900 1 1 2020-01-01T00:00:00 100.00 200.00 12345678900001 12345678900001"}

			require.Error(t, handler.CreateClientsFromLines(ctx, line))
		})

		t.Run("create clients", func(t *testing.T) {
			expectedErr := errors.New("")
			clientService.EXPECT().CreateMulti(ctx, clients).Return(expectedErr)

			require.ErrorIs(t, handler.CreateClientsFromLines(ctx, lines), expectedErr)
		})
	})
}

func Test_client_NormalizeClients(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientService := service.NewMockClientService(ctrl)
	handler := NewClient(clientService)

	t.Run("success", func(t *testing.T) {
		clientService.EXPECT().NormalizeClients(ctx).Return(nil)

		require.NoError(t, handler.NormalizeClients(ctx))
	})

	t.Run("error", func(t *testing.T) {
		clientService.EXPECT().NormalizeClients(ctx).Return(errors.New(""))

		require.Error(t, handler.NormalizeClients(ctx))
	})
}
