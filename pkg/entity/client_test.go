package entity

import (
	"testing"

	"github.com/stretchr/testify/require"

	"neoway/pkg/lib/pointer"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		data    []string
		want    require.ValueAssertionFunc
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "success",
			data:    []string{"12345678900", "1", "1", "2020-01-01", "100.00", "200.00", "12345678900001", "12345678900001"},
			want:    require.NotEmpty,
			wantErr: require.NoError,
		},
		{
			name:    "error/invalid last order date",
			data:    []string{"12345678900", "1", "1", "2020-01-01 00:00:00", "100.00", "200.00", "12345678900001", "12345678900001"},
			want:    require.Nil,
			wantErr: require.Error,
		},
		{
			name:    "error/invalid average ticket",
			data:    []string{"12345678900", "1", "1", "2020-01-01", "100.0a", "200.00", "12345678900001", "12345678900001"},
			want:    require.Nil,
			wantErr: require.Error,
		},
		{
			name:    "error/invalid last order ticket",
			data:    []string{"12345678900", "1", "1", "2020-01-01", "100.00", "200.0a", "12345678900001", "12345678900001"},
			want:    require.Nil,
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.data)
			tt.wantErr(t, err)
			tt.want(t, client)
		})
	}
}

func TestClient_ValidateDocuments(t *testing.T) {
	tests := []struct {
		name                       string
		client                     *Client
		validCPFWant               require.BoolAssertionFunc
		validMostFrequentStoreWant require.BoolAssertionFunc
		validLastOrderStoreWant    require.BoolAssertionFunc
	}{
		{
			name: "all valid",
			client: &Client{
				CPF:               "37078130022",
				MostFrequentStore: pointer.New("11444777000161"),
				LastOrderStore:    pointer.New("11444777000161"),
			},
			validCPFWant:               require.True,
			validMostFrequentStoreWant: require.True,
			validLastOrderStoreWant:    require.True,
		},
		{
			name: "invalid cpf",
			client: &Client{
				CPF:               "37078130021",
				MostFrequentStore: pointer.New("11444777000161"),
				LastOrderStore:    pointer.New("11444777000161"),
			},
			validCPFWant:               require.False,
			validMostFrequentStoreWant: require.True,
			validLastOrderStoreWant:    require.True,
		},
		{
			name: "invalid most frequent store",
			client: &Client{
				CPF:               "37078130022",
				MostFrequentStore: pointer.New("11444777000162"),
				LastOrderStore:    pointer.New("11444777000161"),
			},
			validCPFWant:               require.True,
			validMostFrequentStoreWant: require.False,
			validLastOrderStoreWant:    require.True,
		},
		{
			name: "invalid last order store",
			client: &Client{
				CPF:               "37078130022",
				MostFrequentStore: pointer.New("11444777000161"),
				LastOrderStore:    pointer.New("11444777000162"),
			},
			validCPFWant:               require.True,
			validMostFrequentStoreWant: require.True,
			validLastOrderStoreWant:    require.False,
		},
		{
			name: "all invalid",
			client: &Client{
				CPF:               "37078130021",
				MostFrequentStore: pointer.New("11444777000162"),
				LastOrderStore:    pointer.New("11444777000162"),
			},
			validCPFWant:               require.False,
			validMostFrequentStoreWant: require.False,
			validLastOrderStoreWant:    require.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.client.ValidateDocuments()
			tt.validCPFWant(t, tt.client.ValidCPF)
			tt.validMostFrequentStoreWant(t, tt.client.ValidMostFrequentStore)
			tt.validLastOrderStoreWant(t, tt.client.ValidLastOrderStore)
		})
	}
}
