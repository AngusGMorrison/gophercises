package urlshort

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	store, teardown := setUpStore()
	defer teardown()

	tests := []struct {
		descr, path, wantURL string
		mux                  *http.ServeMux
		wantStatus           int
	}{
		{
			descr:      "redirects when a matching URL is found",
			path:       "/godoc", // matches row in seed file
			wantURL:    "https://godoc.org/",
			mux:        http.DefaultServeMux,
			wantStatus: http.StatusFound,
		},
		{
			descr:      "invokes the fallback handler if path is not matched by DB",
			path:       "/fallback", // does not match row in seed file
			wantURL:    "",
			mux:        http.NewServeMux(), // prevent matching of "/"
			wantStatus: http.StatusOK,
		},
		{
			descr:      "returns 404 Not Found if path is not matched by DB and fallback",
			path:       "/invalid", // does not match seed file or fallback
			wantURL:    "",
			mux:        http.NewServeMux(), // prevent matching of "/"
			wantStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.descr, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			rec := httptest.NewRecorder()
			if test.path != "/invalid" {
				test.mux.HandleFunc(test.path, func(w http.ResponseWriter, r *http.Request) {})
			}
			handler := Handler(store, test.mux)
			handler(rec, req)
			resp := rec.Result()

			if resp.StatusCode != test.wantStatus {
				t.Fatalf(
					"got response status %d %s, want %d %s",
					resp.StatusCode, http.StatusText(resp.StatusCode),
					test.wantStatus, http.StatusText(test.wantStatus),
				)
			}

			if test.wantURL != "" {
				finalURL := resp.Header.Get("Location")
				if finalURL != test.wantURL {
					t.Fatalf("redirected to %s, want %s", finalURL, test.wantURL)
				}
			}
		})
	}
}

func setUpStore() (store *RedirectStore, teardown func()) {
	path := tempfile()
	store, err := OpenRedirectStore(path)
	if err != nil {
		os.Remove(path)
		panic(err)
	}
	err = store.Seed(seedFixture)
	if err != nil {
		store.DB.Close()
		os.Remove(path)
		panic(err)
	}
	teardown = func() {
		store.DB.Close()
		os.Remove(path)
	}
	return
}
