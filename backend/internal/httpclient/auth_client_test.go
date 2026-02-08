package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthHttpClient_GetAuthWellKnownConfig(t *testing.T) {
	t.Run("1. 正常系のテスト", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"issuer": "test"}`))
		}))
		defer server.Close()

		client := NewAuthHttpClient(server.URL)
		resp, err := client.GetAuthWellKnownConfig(t.Context())

		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("2. コンテキストキャンセルのテスト", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewAuthHttpClient(server.URL)
		ctx, cancel := context.WithCancel(t.Context())
		cancel() // 実行前にキャンセル

		resp, err := client.GetAuthWellKnownConfig(ctx)

		require.Error(t, err)
		require.Contains(t, err.Error(), context.Canceled.Error())
		require.Nil(t, resp)
	})

	t.Run("3. タイムアウトのテスト", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // クライアントのタイムアウトより長く待つ
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewAuthHttpClient(server.URL)
		// デフォルトの10秒だと長いので、テスト用に短いタイムアウトのクライアントを上書き
		client.httpClient.Timeout = 50 * time.Millisecond

		resp, err := client.GetAuthWellKnownConfig(t.Context())

		require.Error(t, err)
		require.True(t, assert.Contains(t, err.Error(), "deadline exceeded") || assert.Contains(t, err.Error(), "timeout"))
		require.Nil(t, resp)
	})

	t.Run("4. 不正なURLのテスト", func(t *testing.T) {
		// 制御文字を含むURLは http.NewRequestWithContext でエラーになる
		client := NewAuthHttpClient("http://invalid-url\x7f")
		resp, err := client.GetAuthWellKnownConfig(t.Context())

		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("5. サーバーエラーのテスト", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewAuthHttpClient(server.URL)
		resp, err := client.GetAuthWellKnownConfig(t.Context())

		require.NoError(t, err) // HTTPエラーでもDo自体は成功する
		defer resp.Body.Close()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
