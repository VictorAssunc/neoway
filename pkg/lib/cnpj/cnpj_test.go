package cnpj

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cnpj string
		want require.BoolAssertionFunc
	}{
		{
			name: "valid cnpj",
			cnpj: "11444777000161",
			want: require.True,
		},
		{
			name: "valid cnpj/first digit zero",
			cnpj: "20775878000106",
			want: require.True,
		},
		{
			name: "valid cnpj/second digit zero",
			cnpj: "91057935000160",
			want: require.True,
		},
		{
			name: "invalid cnpj/invalid length",
			cnpj: "1144477700016",
			want: require.False,
		},
		{
			name: "invalid cnpj/non-digit",
			cnpj: "1144477700016a",
			want: require.False,
		},
		{
			name: "invalid cnpj/invalid first digit",
			cnpj: "11444777000171",
			want: require.False,
		},
		{
			name: "invalid cnpj/invalid second digit",
			cnpj: "11444777000162",
			want: require.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want(t, Validate(tt.cnpj))
		})
	}
}
