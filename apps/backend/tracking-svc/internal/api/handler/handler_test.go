package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondSuccess(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := map[string]string{
		"status": "ok",
	}

	RespondSuccess(recorder, http.StatusOK, data)

	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected content-type application/json, got %s", ct)
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}


