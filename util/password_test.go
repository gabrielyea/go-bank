package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	input := "s3cret"
	res, err := HashPassword(input)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	err = IsValidPassword(input, res)
	require.NoError(t, err)

	err = IsValidPassword("hola", res)
	require.Error(t, err)
}
