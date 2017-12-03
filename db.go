package main

import (
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
	"log"
	"encoding/binary"
)

var wordsBucket = "words"

type DBConnect struct {
	path string
}

func NewDB(path string) *DBConnect {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(wordsBucket))
		return err
	})
	if err != nil {
		panic(err)
	}

	return &DBConnect{path}
}

func (c *DBConnect) AddWord(word Word) (*Word, error) {
	db, err := open(c)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(wordsBucket))

		// Marshal user data into bytes.
		buf, err := json.Marshal(&word)
		if err != nil {
			return err
		}

		// Persist bytes to environment bucket.
		return b.Put([]byte(word.Text), buf)
	})
	if err != nil {
		log.Printf("Error on add word to DB. %s\n", err)
		return nil, err
	}

	return &word, nil

}

func (c *DBConnect) CountWords() (int, error) {
	db, err := open(c)
	if err != nil {
		return 0, err
	}

	defer db.Close()

	var count int
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(wordsBucket))

		count = b.Stats().KeyN

		return nil
	})
	if err != nil {
		log.Printf("Error on get count words in DB. %s\n", err)
		return 0, err
	}

	return count, nil
}

func (c *DBConnect) GetWords(text string) (*Word, error) {
	db, err := open(c)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var data []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(wordsBucket))
		data = b.Get([]byte(text))
		return nil
	})
	if err != nil {
		return nil, err
	}
	//if not found data by key
	if len(data) == 0 {
		return nil, nil
	}
	//parse byte array
	var env Word
	err = json.Unmarshal(data, &env)
	if err != nil {
		log.Printf("Error on get word from DB. %s\n", err)
		return nil, err
	}

	return &env, nil
}

func open(c *DBConnect) (*bolt.DB, error) {
	db, err := bolt.Open(c.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Printf("Error on open connections with DB. %s\n", err)
		return nil, err
	}

	return db, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
