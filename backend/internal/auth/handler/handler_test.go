package handler

import (
	"echo-sample2/internal/httpclient"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_RegisterAuthRoutes(t *testing.T) {
	t.Run("正常系: GetAuthWellKnownConfigが成功する場合", func(t *testing.T) {
		// 1. AuthHttpClientが呼び出すモックサーバーの作成
		mockAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/.well-known/openid-configuration", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"issuer": "http://mock-auth-server"}`))
		}))
		defer mockAuthServer.Close()

		// 2. AuthHandlerのセットアップ
		authClient := httpclient.NewAuthHttpClient(mockAuthServer.URL + "/.well-known/openid-configuration")
		h := NewAuthHandler(authClient)
		e := echo.New()
		h.RegisterAuthRoutes(e)

		// 3. リクエストの実行
		req := httptest.NewRequest(http.MethodGet, "/auth-well-known-config", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// 4. 検証
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		require.JSONEq(t, `{"issuer": "http://mock-auth-server"}`, rec.Body.String())
	})

	t.Run("異常系: GetAuthWellKnownConfigがエラーを返す場合", func(t *testing.T) {
		// 制御文字を含む不正なURLを指定して、httpclient.GetAuthWellKnownConfig でエラーを発生させる
		authClient := httpclient.NewAuthHttpClient("http://invalid-url\x7f")
		h := NewAuthHandler(authClient)
		e := echo.New()
		h.RegisterAuthRoutes(e)

		req := httptest.NewRequest(http.MethodGet, "/auth-well-known-config", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// ハンドラー内で err が返されるため、Echoのデフォルトエラーハンドラーにより 500 が返される
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("異常系: 認証サーバーがエラーを返す場合", func(t *testing.T) {
		mockAuthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer mockAuthServer.Close()

		authClient := httpclient.NewAuthHttpClient(mockAuthServer.URL)
		h := NewAuthHandler(authClient)
		e := echo.New()
		h.RegisterAuthRoutes(e)

		req := httptest.NewRequest(http.MethodGet, "/auth-well-known-config", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
}
