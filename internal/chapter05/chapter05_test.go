package chapter05_test

import (
	"github.com/nimsaysm/goProgrammingExercises/internal/chapter05"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTMLOutline(t *testing.T) {
	t.Run("Successful getting expected output", func(t *testing.T) {
		content := `
			<!DOCTYPE html>
			<html>
			<head><title>Test</title></head>
			<body>
				<h1>Hello, World!</h1>
				<p>This is a test.</p>
			</body>
			</html>
		`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content))
		}))
		defer server.Close()

		err := chapter05.HTMLOutline([]string{server.URL})
		require.NoError(t, err)
	})
}
