package redisI

import (
	"fmt"
	"time"
)

type Key struct {
	Key     string
	Expired time.Duration
}

var (
	AuthNonce = Key{"auth_nonce:%s", time.Minute * 60 * 24}
)

func (r *Key) SKey(prefix string) string {
	return fmt.Sprintf(r.Key, prefix)
}
