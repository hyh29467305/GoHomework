package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type LocalCodeCache struct {
	cache *lru.Cache
	//普通锁，或者说写锁
	lock sync.Mutex
	//读写锁
	expiration time.Duration
}

func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) *LocalCodeCache {
	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		//说明没有验证码
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expiration),
		})
		return nil
	}
	item, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if item.expire.Sub(now) > time.Minute*9 {
		return ErrCodeSendTooMany
	}
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrKeyNoExist
	}
	item, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}
	if item.cnt <= 0 {
		return false, ErrCodeVefifyTooMany
	}
	item.cnt--
	return item.code == inputCode, nil
}

func (l *LocalCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code string
	//可验证次数
	cnt int
	//过期时间
	expire time.Time
}
