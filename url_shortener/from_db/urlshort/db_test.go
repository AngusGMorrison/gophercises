package urlshort

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

const seedFixture = "fixtures/seeds_test.yaml"

func TestOpenRedirectStore(t *testing.T) {
	oldDBName := dbName
	dbName = tempfile()
	defer os.Remove(dbName)

	// Test DB creation in main routine to ensure failure of test suite and proper teardown if it
	// fails.
	store, err := OpenRedirectStore()
	if err != nil {
		t.Fatalf("opening DB: %v", err)
	}
	defer store.DB.Close()

	t.Run(fmt.Sprintf("creates bucket %q", bucketName), func(t *testing.T) {
		bucketExists := false
		store.DB.View(func(tx *bolt.Tx) error {
			if b := tx.Bucket([]byte(bucketName)); b != nil {
				bucketExists = true
			}
			return nil
		})
		if !bucketExists {
			t.Fatalf("default bucket %q was not created", bucketName)
		}
	})

	dbName = oldDBName
}

func TestSeed(t *testing.T) {
	oldDBName := dbName
	dbName = tempfile()
	defer os.Remove(dbName)

	store, err := OpenRedirectStore()
	if err != nil {
		panic(err)
	}
	defer store.DB.Close()

	// Test for errors
	t.Run("seeds the DB without errors", func(t *testing.T) {
		if err := store.Seed(seedFixture); err != nil {
			t.Fatalf("seeding DB: %v", err)
		}
	})

	redirects, _ := readSeedFile(seedFixture) // ignore error; must succeed if Seed succeeded
	t.Run("adds the correct number of rows to the DB", func(t *testing.T) {
		store.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			keyN := b.Stats().KeyN
			if keyN != len(redirects) {
				t.Fatalf("found %d entries in the DB, want %d", keyN, len(redirects))
			}
			return nil
		})
	})

	t.Run("adds the correct key-value pairs to the DB", func(t *testing.T) {
		r := redirects[0]
		store.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			v := b.Get([]byte(r.Path))
			if gotURL := string(v); gotURL != r.URL {
				t.Fatalf("entry for path %q had URL %q, want %q", r.Path, gotURL, r.URL)
			}
			return nil
		})
	})

	dbName = oldDBName
}

func tempfile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
