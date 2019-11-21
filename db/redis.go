package db

import (
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

func NewRedisConnPool(url string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     100,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
