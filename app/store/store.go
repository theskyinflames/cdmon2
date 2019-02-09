package store

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/theskyinflames/cdmon2/app"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/theskyinflames/cdmon2/app/config"
	"github.com/theskyinflames/cdmon2/app/domain"

	"github.com/go-redis/redis"
)

func init() {
	gob.Register(domain.Server{})
	gob.Register(domain.Hosting{})
}

type (
	Store struct {
		log  *logrus.Logger
		cfg  *config.Config
		conn *redis.Client
	}
)

func NewStore(cfg *config.Config, log *logrus.Logger) (*Store, error) {
	s := &Store{
		log: log,
		cfg: cfg,
	}
	err := s.Connect()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (_ *Store) ItemToJSON(item interface{}) ([]byte, error) {
	return json.Marshal(item)
}

func (_ *Store) FromJSONToItem(b []byte, item interface{}) (interface{}, error) {
	err := json.Unmarshal(b, item)
	return item, err
}

// Item serializing
func (_ *Store) ItemToGob(item interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Item deserializing
func (_ *Store) FromGobToItem(b []byte, item interface{}) (interface{}, error) {
	buff := bytes.Buffer{}
	buff.Write(b)
	d := gob.NewDecoder(&buff)
	err := d.Decode(item) // item must be a pointer
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Store) Connect() error {

	s.log.Infof("redis connection at %s", s.cfg.RedisAddr)
	s.conn = redis.NewClient(&redis.Options{
		Addr:     s.cfg.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := s.conn.Ping()
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (s *Store) Flush() error {
	res := s.conn.FlushAll()
	return res.Err()
}

func (s *Store) Close() error {
	return s.conn.Close()
}

func (s *Store) Get(key string, item interface{}) (interface{}, error) {
	res := s.conn.Get(key)
	err := res.Err()
	if err != nil {
		switch errors.Cause(err) {
		case redis.Nil:
			return nil, app.DbErrorNotFound
		default:
			return nil, err
		}
	}

	bin, err := res.Result()
	if err != nil {
		return nil, err
	}
	return s.FromGobToItem([]byte(bin), item)
}

func (s *Store) GetAll(pattern string, emptyRecordFunc config.EmptyRecordFunc) ([]interface{}, error) {
	if len(pattern) == 0 {
		pattern = "*"
	}
	res := s.conn.Keys(pattern)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var keys []string = res.Val()
	s.log.Infof("retrieved %d hostings to GetAll()", len(keys))
	slice := make([]interface{}, len(keys))
	for z, k := range keys {
		item, err := s.Get(k, emptyRecordFunc())
		if err != nil {
			return nil, err
		}
		slice[z] = item
	}
	return slice, nil
}

func (s *Store) Set(key string, item interface{}) error {
	bin, err := s.ItemToGob(item)
	if err != nil {
		return nil
	}
	return s.conn.Set(key, bin, 0).Err()
}

func (s *Store) Remove(key string) error {
	return s.conn.Del(key).Err()
}
