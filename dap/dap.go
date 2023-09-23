package dap

import (
	"errors"
	"fmt"

	"github.com/raja-dettex/dap/schema"
	"go.etcd.io/bbolt"
)

var (
	ext string = "dap"
)

type Dap struct {
	opts *Options
	Db   *bbolt.DB
}

type Collection struct {
	Bucket *bbolt.Bucket
}

type Any struct {
}

// factory to open a db

func New(dbName string) (*Dap, error) {
	name := fmt.Sprintf("%s.%s", dbName, ext)
	db, err := bbolt.Open(name, 0600, nil)
	if err != nil {
		return nil, err
	}
	options := &Options{
		Encoder: &JSONEncoder{},
		Decoder: &JSONDecoder{},
		DbName:  name,
	}
	return &Dap{Db: db, opts: options}, nil
}

func (d *Dap) CreateCollection(collName string) (*Collection, error) {
	var collection = &Collection{}
	tx, err := d.Db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	bucket, err := tx.CreateBucketIfNotExists([]byte(collName))
	if err != nil {
		return nil, err
	}
	collection.Bucket = bucket
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (d *Dap) Insert(collName, key string, data schema.Data) error {
	tx, err := d.Db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	bucket, err := tx.CreateBucketIfNotExists([]byte(collName))
	if err != nil {
		return err
	}
	val, err := d.opts.Encoder.Encode(data)
	if err != nil {
		return err
	}
	err = bucket.Put([]byte(key), val)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *Dap) Select(collName, key string) (schema.Data, error) {
	tx, err := d.Db.Begin(true)
	if err != nil {
		return nil, err
	}
	bucket := tx.Bucket([]byte(collName))
	if bucket == nil {
		return nil, errors.New("collection does not exist")
	}
	val := bucket.Get([]byte(key))
	data := make(schema.Data)
	decoded, err := d.opts.Decoder.Decode(val, data)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return decoded, nil
}
