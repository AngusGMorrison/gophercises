package urlshort

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// RedirectStore wraps a boltdb instance, enabling attachment of methods such as Seed.
type RedirectStore struct {
	DB *bolt.DB
}

var (
	dbName     = "redirects.db"
	bucketName = "redirects"
)

// OpenRedirectStore returns a new boltdb instance as a RedirectDB with a bucket to store redirects.
// The caller is responsible for closing the DB when done.
func OpenRedirectStore() (*RedirectStore, error) {
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	}); err != nil {
		db.Close()
		return nil, err
	}

	return &RedirectStore{db}, nil
}

// Seed parses a YAML file and adds all path:url pairs found to rs' underlying database.
func (rs *RedirectStore) Seed(path string) error {
	redirects, err := readSeedFile(path)
	if err != nil {
		return err
	}

	err = rs.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		for _, r := range redirects {
			err := b.Put([]byte(r.Path), []byte(r.URL))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func readSeedFile(path string) ([]Redirect, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var redirects []Redirect
	if err = yaml.Unmarshal(data, &redirects); err != nil {
		return nil, err
	}
	return redirects, nil
}
