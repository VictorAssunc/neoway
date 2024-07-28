package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"neoway/pkg/entity"
	"neoway/pkg/lib/pointer"
	repository "neoway/pkg/repository/mock"
)

func TestNewClient(t *testing.T) {
	require.NotNil(t, NewClient(nil))
}

func TestClient_CreateMulti(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := repository.NewMockClientRepository(ctrl)
	service := NewClient(repositoryMock)

	clients := []*entity.Client{
		{CPF: "12345678901"},
		{CPF: "12345678902"},
	}

	t.Run("success", func(t *testing.T) {
		repositoryMock.EXPECT().CreateMulti(ctx, clients).Return(nil)

		require.NoError(t, service.CreateMulti(ctx, clients))
	})

	t.Run("error", func(t *testing.T) {
		repositoryMock.EXPECT().CreateMulti(ctx, clients).Return(errors.New(""))

		require.Error(t, service.CreateMulti(ctx, clients))
	})
}

func TestClient_NormalizeClients(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repositoryMock := repository.NewMockClientRepository(ctrl)
	service := NewClient(repositoryMock)

	clients := []*entity.Client{
		{
			CPF:                    "37078130022",
			MostFrequentStore:      pointer.New("11444777000161"),
			LastOrderStore:         pointer.New("11444777000161"),
			ValidCPF:               true,
			ValidMostFrequentStore: true,
			ValidLastOrderStore:    true,
		},
		{
			CPF:                    "37078130021",
			MostFrequentStore:      pointer.New("11444777000162"),
			LastOrderStore:         pointer.New("11444777000162"),
			ValidCPF:               true,
			ValidMostFrequentStore: true,
			ValidLastOrderStore:    true,
		},
	}

	t.Run("success", func(t *testing.T) {
		repositoryMock.EXPECT().GetPaginated(ctx, SearchLimit, 0).Return(clients, nil)
		repositoryMock.EXPECT().UpdateMultiDocumentsValidity(ctx, clients).Return(nil)
		repositoryMock.EXPECT().GetPaginated(ctx, SearchLimit, 1000).Return(nil, nil)

		require.NoError(t, service.NormalizeClients(ctx))
	})

	t.Run("error", func(t *testing.T) {
		t.Run("get client", func(t *testing.T) {
			expectedErr := errors.New("")
			repositoryMock.EXPECT().GetPaginated(ctx, SearchLimit, 0).Return(nil, expectedErr)

			require.ErrorIs(t, service.NormalizeClients(ctx), expectedErr)
		})

		t.Run("update clients", func(t *testing.T) {
			expectedErr := errors.New("")
			repositoryMock.EXPECT().GetPaginated(ctx, SearchLimit, 0).Return(clients, nil)
			repositoryMock.EXPECT().UpdateMultiDocumentsValidity(ctx, clients).Return(expectedErr)

			require.ErrorIs(t, service.NormalizeClients(ctx), expectedErr)
		})
	})
}
