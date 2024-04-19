package redis

import (
	"fmt"
	"time"
)

type RedisKey struct {
	Key     string
	Expired time.Duration
}

var (
	AUTH_NONCE = RedisKey{"auth_nonce:%s", time.Minute * 60 * 24}
)

func (r *RedisKey) SKey(prefix string) string {
	return fmt.Sprintf(r.Key, prefix)
}
