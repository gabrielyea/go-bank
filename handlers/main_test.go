package handlers

import (
	"os"
	"testing"
	"time"

	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, repo repo.Store) *Server {
	config := util.Config{
		SymmetricKey:  util.RandomString(32),
		TokenDuration: time.Minute,
	}

	h := NewHandler(repo)

	server := SetUpServer(config, h)
	require.NotNil(t, server)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
