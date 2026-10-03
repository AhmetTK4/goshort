package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AhmetTK4/goshort/storage"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func testAPI(t *testing.T) http.Handler {
	t.Helper()
	server := miniredis.RunT(t)
	previous := storage.RDB
	storage.RDB = redis.NewClient(&redis.Options{Addr: server.Addr()})
	client := storage.RDB
	t.Cleanup(func() { client.Close(); storage.RDB = previous })
	t.Setenv("BASE_URL", "https://short.example")
	return newRouter()
}

func TestShortenRedirectAndCount(t *testing.T) {
	api := testAPI(t)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(`{"url":"https://example.com/page"}`))
	request.Header.Set("Content-Type", "application/json")
	api.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("shorten: %d %s", response.Code, response.Body)
	}
	var created struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(created.ShortURL, "https://short.example/g/") {
		t.Fatalf("unexpected URL: %s", created.ShortURL)
	}
	code := strings.TrimPrefix(created.ShortURL, "https://short.example/g/")
	redirected := httptest.NewRecorder()
	api.ServeHTTP(redirected, httptest.NewRequest(http.MethodGet, "/g/"+code, nil))
	if redirected.Code != http.StatusFound || redirected.Header().Get("Location") != "https://example.com/page" {
		t.Fatalf("redirect: %d %s", redirected.Code, redirected.Header().Get("Location"))
	}
	stats := httptest.NewRecorder()
	api.ServeHTTP(stats, httptest.NewRequest(http.MethodGet, "/api/stats/"+code, nil))
	var counted struct {
		Clicks string `json:"clicks"`
	}
	if err := json.Unmarshal(stats.Body.Bytes(), &counted); err != nil {
		t.Fatal(err)
	}
	if counted.Clicks != "1" {
		t.Fatalf("clicks = %s, want 1", counted.Clicks)
	}
}

func TestRejectsInvalidURLs(t *testing.T) {
	api := testAPI(t)
	for _, input := range []string{"", "javascript:alert(1)", "/relative", "ftp://example.com", "https://user:pass@example.com"} {
		t.Run(input, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{"url": input})
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			api.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d", response.Code)
			}
		})
	}
}

func TestUnknownShortCodeReturnsNotFound(t *testing.T) {
	api := testAPI(t)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/g/missing", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestCustomShortCode(t *testing.T) {
	api := testAPI(t)
	body := `{"url":"https://example.com/custom","custom_code":"mycode"}`
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("custom code: %d %s", response.Code, response.Body)
	}
	var created struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	expected := "https://short.example/g/mycode"
	if created.ShortURL != expected {
		t.Fatalf("got %s, want %s", created.ShortURL, expected)
	}
}

func TestCustomCodeValidation(t *testing.T) {
	api := testAPI(t)
	tests := []struct {
		name string
		code string
		want int
	}{
		{"too short", "ab", http.StatusBadRequest},
		{"too long", "abcdefghij1234567890x", http.StatusBadRequest},
		{"invalid chars", "my-code", http.StatusBadRequest},
		{"valid", "valid123", http.StatusOK},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]string{"url": "https://example.com", "custom_code": tc.code})
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			api.ServeHTTP(response, request)
			if response.Code != tc.want {
				t.Fatalf("status = %d, want %d", response.Code, tc.want)
			}
		})
	}
}

func TestCustomCodeConflict(t *testing.T) {
	api := testAPI(t)
	body := `{"url":"https://example.com/first","custom_code":"conflict"}`
	request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("first request: %d %s", response.Code, response.Body)
	}
	
	body2 := `{"url":"https://example.com/second","custom_code":"conflict"}`
	request2 := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(body2))
	request2.Header.Set("Content-Type", "application/json")
	response2 := httptest.NewRecorder()
	api.ServeHTTP(response2, request2)
	if response2.Code != http.StatusConflict {
		t.Fatalf("second request: got %d, want %d", response2.Code, http.StatusConflict)
	}
}
