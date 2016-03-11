package opendmm

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

func NewDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MovieMeta"))
		if err != nil {
			return fmt.Errorf("Error creating bucket MovieMeta: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func writeMetaToDB(meta MovieMeta, db *bolt.DB) error {
	bdata, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MovieMeta"))
		return bucket.Put([]byte(meta.Code), bdata)
	})
}
