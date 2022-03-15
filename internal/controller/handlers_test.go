package controller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/repository"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SendTestRequest(t *testing.T, ts *httptest.Server, method, path, contentType string, body io.Reader) (*http.Response, string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()
	return resp, string(respBody)
}

type want struct {
	responseStatusCode  int
	responseParams      map[string]string
	responseContentType string
	responseBody        string
}

type request struct {
	url, method, contentType, body string
}

const testAddr = "http://localhost:8080/"

func TestUserRegister(t *testing.T) {
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "test #1. Negative. POST with empty body",
			request: request{
				url:    "/api/user/register",
				method: http.MethodPost,
				body:   "",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "",
			},
		},
		{
			name: "test #2. Negative. POST with empty password",
			request: request{
				url:    "/api/user/register",
				method: http.MethodPost,
				body:   "{\n  \"login\": \"\",\n  \"password\": \"string\"\n}",
			},
			want: want{
				responseStatusCode: http.StatusBadRequest,
				responseParams:     nil,
				responseBody:       "",
			},
		}}
	//urMock := &mocks.UserRepo{}
	//urMock.On("AddUser" ,mock.Anything).Return(errors.New("pgerrcode.UniqueViolation"))
	service := gophermart.NewGophermart(testAddr, repository.NewUserRepository(nil), repository.NewUserAccountRepository(nil), repository.NewUserOrderRepository(nil), nil)

	r := NewRouter(context.Background(), service)
	ts := httptest.NewServer(r)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := SendTestRequest(t, ts, tt.request.method, tt.request.url, "", nil)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.responseStatusCode, resp.StatusCode)

		})
	}
}
