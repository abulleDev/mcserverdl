package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testJSONStruct struct {
	Array []struct {
		Number  int    `json:"number"`
		String  string `json:"string"`
		Boolean bool   `json:"boolean"`
	} `json:"array"`
}

var testJSONDataConst = testJSONStruct{
	Array: []struct {
		Number  int    `json:"number"`
		String  string `json:"string"`
		Boolean bool   `json:"boolean"`
	}{
		{
			Number:  255,
			String:  "text",
			Boolean: true,
		},
	},
}

func TestFetchJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(testJSONDataConst); err != nil {
				t.Fatalf("failed to encode test JSON: %v", err)
			}
		}))
		defer testServer.Close()

		var testJSONData testJSONStruct
		if err := FetchJSON(testServer.URL, &testJSONData); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(testJSONData.Array) != 1 {
			t.Fatalf("expected array length 1, got: %d", len(testJSONData.Array))
		}
		if len(testJSONData.Array) != len(testJSONDataConst.Array) ||
			testJSONData.Array[0].Number != testJSONDataConst.Array[0].Number ||
			testJSONData.Array[0].String != testJSONDataConst.Array[0].String ||
			testJSONData.Array[0].Boolean != testJSONDataConst.Array[0].Boolean {
			t.Errorf("unexpected struct values: %+v", testJSONData.Array[0])
		}
	})

	t.Run("not found", func(t *testing.T) {
		testServer := httptest.NewServer(http.NotFoundHandler())
		defer testServer.Close()

		var testJSONData testJSONStruct
		if err := FetchJSON(testServer.URL, &testJSONData); err == nil {
			t.Error("expected error for 404 response, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "not a json")
		}))
		defer testServer.Close()

		var testJSONData testJSONStruct
		if err := FetchJSON(testServer.URL, &testJSONData); err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})

	t.Run("redirect followed", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/redirect" {
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(testJSONDataConst); err != nil {
				t.Fatalf("failed to encode test JSON: %v", err)
			}
		}))
		defer testServer.Close()

		var testJSONData testJSONStruct
		if err := FetchJSON(testServer.URL+"/redirect", &testJSONData); err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}
