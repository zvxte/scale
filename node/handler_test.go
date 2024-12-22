package node

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIndex(t *testing.T) {
	logger := log.Default()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler := getIndex(logger)
	handler(recorder, req)

	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error(resp.StatusCode, resp.Status)
	}
}
