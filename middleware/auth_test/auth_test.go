package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gabriel/gabrielyea/go-bank/handlers"
	"github.com/gabriel/gabrielyea/go-bank/middleware"
	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gabriel/gabrielyea/go-bank/token"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, repo repo.Store) *handlers.Server {
	config := util.Config{
		SymmetricKey:  util.RandomString(32),
		TokenDuration: time.Minute,
	}

	h := handlers.NewHandler(repo)

	server := handlers.SetUpServer(config, h)
	require.NotNil(t, server)
	return server
}

func AddAuthorization(
	t *testing.T,
	request *http.Request,
	tm token.Maker,
	authType string,
	username string,
	duration time.Duration,
) {
	token, err := tm.CreateToken(username, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(middleware.AuthKeys()["headerKey"], authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tMaker token.Maker) {
				AddAuthorization(t, request, tMaker, middleware.AuthKeys()["auth"], "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "NoAuth",
			setupAuth: func(t *testing.T, request *http.Request, tMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "UnsupportedAuth",
			setupAuth: func(t *testing.T, request *http.Request, tMaker token.Maker) {
				AddAuthorization(t, request, tMaker, "notbearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "TokenExpired",
			setupAuth: func(t *testing.T, request *http.Request, tMaker token.Maker) {
				AddAuthorization(t, request, tMaker, middleware.AuthKeys()["auth"], "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)
			authPath := "/auth"
			server.Router.GET(authPath, middleware.AuthMiddleware(server.TokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.TokenMaker)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)

		})
	}
}
