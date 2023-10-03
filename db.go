package main

import (
	"go.etcd.io/bbolt"
)

func setupDatabase() (*bbolt.DB, error) {
	db, err := bbolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("secrets"))
		return err
	})

	return db, err
}
