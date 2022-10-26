package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_repo "github.com/gabriel/gabrielyea/go-bank/db/mock"
	"github.com/gabriel/gabrielyea/go-bank/repo"

	"github.com/gabriel/gabrielyea/go-bank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountHandler(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(r *mock_repo.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(r *mock_repo.MockStore) {
				r.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rr.Code)
				requireBodyMAtchAccount(t, rr.Body, account)

			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(r *mock_repo.MockStore) {
				r.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repo.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rr.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(r *mock_repo.MockStore) {
				r.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(repo.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rr.Code)
			},
		},
		{
			name:      "InvalidId",
			accountID: 0,
			buildStubs: func(r *mock_repo.MockStore) {
				r.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rr.Code)
			},
		},
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

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			server.Router.ServeHTTP(recorder, request)
			require.NoError(t, err)

			tc.checkResponse(t, recorder)
		})
	}

}

func requireBodyMAtchAccount(t *testing.T, body *bytes.Buffer, account repo.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var testAccount repo.Account
	err = json.Unmarshal(data, &testAccount)

	require.NoError(t, err)
	require.Equal(t, account, testAccount)

}

func randomAccount() repo.Account {
	return repo.Account{
		ID:       util.RandomInt(1, 10000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
