package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mock_repo "github.com/gabriel/gabrielyea/go-bank/db/mock"
	"github.com/gabriel/gabrielyea/go-bank/repo"
	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      repo.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(repo.CreateUserParams)
	if !ok {
		return false
	}

	err := util.IsValidPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg repo.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_repo.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.UserName,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mock_repo.MockStore) {
				arg := repo.CreateUserParams{
					UserName: user.UserName,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		// {
		// 	name: "InternalError",
		// 	body: gin.H{
		// 		"username":  user.UserName,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mock_repo.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(repo.User{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "DuplicateUsername",
		// 	body: gin.H{
		// 		"user_name": user.UserName,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mock_repo.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(1).
		// 			Return(repo.User{}, &pq.Error{Code: "23505"})
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusForbidden, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidUsername",
		// 	body: gin.H{
		// 		"user_name": "invalid-user#1",
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mock_repo.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "InvalidEmail",
		// 	body: gin.H{
		// 		"user_name": user.UserName,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     "invalid-email",
		// 	},
		// 	buildStubs: func(store *mock_repo.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
		// {
		// 	name: "TooShortPassword",
		// 	body: gin.H{
		// 		"user_name": user.UserName,
		// 		"password":  "123",
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mock_repo.MockStore) {
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), gomock.Any()).
		// 			Times(0)
		// 	},
		// 	checkResponse: func(recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repo.NewMockStore(ctrl)
			tc.buildStubs(repo)
			server := NewTestServer(t, repo)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// func TestLoginUserAPI(t *testing.T) {
// 	user, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		buildStubs    func(store *mock_repo.MockStore)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": user.UserName,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mock_repo.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.UserName)).
// 					Times(1).
// 					Return(user, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "UserNotFound",
// 			body: gin.H{
// 				"username": "NotFound",
// 				"password": password,
// 			},
// 			buildStubs: func(store *mock_repo.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(repo.User{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "IncorrectPassword",
// 			body: gin.H{
// 				"username": user.UserName,
// 				"password": "incorrect",
// 			},
// 			buildStubs: func(store *mock_repo.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(user.UserName)).
// 					Times(1).
// 					Return(user, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"username": user.UserName,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mock_repo.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(repo.User{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidUsername",
// 			body: gin.H{
// 				"username": "invalid-user#1",
// 				"password": password,
// 			},
// 			buildStubs: func(store *mock_repo.MockStore) {
// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			repo := mock_repo.NewMockStore(ctrl)
// 			tc.buildStubs(repo)
// 			h := NewHandler(repo)
// 			server := SetUpServer(h)

// 			recorder := httptest.NewRecorder()

// 			// Marshal body data to JSON
// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/login"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.Router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

func randomUser(t *testing.T) (user repo.User, password string) {
	password = util.RandomOwner()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = repo.User{
		UserName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user repo.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser repo.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.UserName, gotUser.UserName)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
