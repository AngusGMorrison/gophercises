package urlshort

import (
	"net/http"

	"github.com/boltdb/bolt"
)

// A Redirect maps a short path to the URL it redirects to.
type Redirect struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// Handler queries db for a URL matching the request's path and redirects to the URL if found. If
// no match exists, the fallback is invoked.
func Handler(store *RedirectStore, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		var foundURL string
		store.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			v := b.Get([]byte(path))
			if v != nil {
				foundURL = string(v)
			}
			return nil
		})

		if len(foundURL) != 0 {
			http.Redirect(w, r, string(foundURL), http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}
