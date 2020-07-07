package urlshort

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

const seedFixture = "fixtures/seeds_test.yaml"

func TestOpenRedirectStore(t *testing.T) {
	oldDBName := dbName
	dbName = "test.db"

	// Test RedirectStore is created
	store, err := OpenRedirectStore()
	if err != nil {
		t.Errorf("opening DB: %v", err)
	}

	// Test bucket is created
	bucketExists := false
	store.DB.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(bucketName)); b != nil {
			bucketExists = true
		}
		return nil
	})
	if !bucketExists {
		t.Errorf("default bucket %q was not created", bucketName)
		store.DB.Close()
	}

	// Clean up
	store.DB.Close()
	os.Remove(dbName)
	dbName = oldDBName
}

func TestSeed(t *testing.T) {
	oldDBName := dbName
	dbName = "test.db"
	store, err := OpenRedirectStore()
	if err != nil {
		t.Errorf("failed to create test DB: %v", err)
	}
	defer store.DB.Close()

	// Test for errors
	if err := store.Seed(seedFixture); err != nil {
		t.Errorf("seeding DB: %v", err)
	}

	// Test key count
	redirects, _ := readSeedFile(seedFixture) // ignore error; must succeed if Seed succeeded
	store.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		keyN := b.Stats().KeyN
		if keyN != len(redirects) {
			t.Errorf("found %d entries in the DB, want %d", keyN, len(redirects))
		}
		return nil
	})

	// Test key-value pair integrity
	r := redirects[0]
	store.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(r.Path))
		if gotURL := string(v); gotURL != r.URL {
			t.Errorf("entry for path %q had URL %q, want %q", r.Path, gotURL, r.URL)
		}
		return nil
	})

	// Test false negatives
	invalidPath := "/invalid"
	store.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(invalidPath))
		if v != nil {
			t.Errorf("unexpected URL %q for invalid path %q, want nil", v, invalidPath)
		}
		return nil
	})

	store.DB.Close()
	os.Remove(dbName)
	dbName = oldDBName
}
