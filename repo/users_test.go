package repo

import (
	"context"
	"testing"

	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashed, err := util.HashPassword(util.RandomOwner())
	arg := CreateUserParams{
		UserName:       util.RandomOwner(),
		HashedPassword: hashed,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.UserName, arg.UserName)
	require.Equal(t, user.Email, arg.Email)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
}
