package cpf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cpf  string
		want require.BoolAssertionFunc
	}{
		{
			name: "valid cpf",
			cpf:  "37078130022",
			want: require.True,
		},
		{
			name: "valid cpf/first digit zero",
			cpf:  "05504106001",
			want: require.True,
		},
		{
			name: "valid cpf/second digit zero",
			cpf:  "15976239030",
			want: require.True,
		},
		{
			name: "invalid cpf/same digits",
			cpf:  "11111111111",
			want: require.False,
		},
		{
			name: "invalid cpf/invalid length",
			cpf:  "3707813002",
			want: require.False,
		},
		{
			name: "invalid cpf/non-digit",
			cpf:  "3707813002a",
			want: require.False,
		},
		{
			name: "invalid cpf/invalid first digit",
			cpf:  "37078130012",
			want: require.False,
		},
		{
			name: "invalid cpf/invalid second digit",
			cpf:  "37078130021",
			want: require.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want(t, Validate(tt.cpf))
		})
	}
}
