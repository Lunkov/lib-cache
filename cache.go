package cache

import (
  "time"
  "reflect"
  "github.com/gomodule/redigo/redis"
  "github.com/golang/glog"
  "sync"
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type ItemStr struct {
	Object     string
	Expiration int64
}

type Cache struct {
	*cache
	// If this is confusing, see the comment at the bottom of New()
}

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
  
  modeUndef  = 0
  modeMemory = 1
  modeRedis  = 2
)

type cache struct {
  mode                int
	defaultExpiration   int64
	items               map[string]Item
  itemsStr            map[string]ItemStr
	mu                  sync.RWMutex
  redisMaxConnections int
  redisURL            string
  redisPool           *redis.Pool
}

func (c *cache) HasError() bool {
  return c.mode == modeUndef
}

func (c *cache) Mode() string {
  switch(c.mode) {
    case modeMemory:
      return "memory"
    case modeRedis:
      return "redis"
  }
  return "undefined"
}

func (c *cache) Count() int64 {
  switch(c.mode) {
    case modeMemory:
      return c.memCount()
    case modeRedis:
      return c.redisCount()
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
  return -1
}

func (c *cache) Clear() {
  switch(c.mode) {
    case modeMemory:
      c.memClear()
      break
    case modeRedis:
      c.redisClear()
      break
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
}

func (c *cache) DefaultExpiration() int64 {
  return c.defaultExpiration
}

func (c *cache) Set(key string, x interface{}) {
  switch(c.mode) {
    case modeMemory:
      c.memSet(key, x)
      break
    case modeRedis:
      c.redisSet(key, x)
      break
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
}

func (c *cache) GetType(myvar interface{}) string {
  t := reflect.TypeOf(myvar)
  if t == nil {
    return "<nil>"
  }
  if t.Kind() == reflect.Ptr {
    return "*" + t.Elem().Name()
  } else {
    return t.Name()
  }
}

func (c *cache) Get(key string, x interface{}) (interface{}, bool) {
  switch(c.mode) {
    case modeMemory:
      return c.memGet(key, x)
    case modeRedis:
      return c.redisGet(key, x)
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
  return nil, false
}

func (c *cache) GetAll2JSON(x interface{}) []byte {
  switch(c.mode) {
    case modeMemory:
      return c.memGetAll2JSON(x)
    case modeRedis:
      return c.redisGetAll2JSON(x)
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
  return []byte("[]")
}

func (c *cache) SetStr(key string, x string) {
  switch(c.mode) {
    case modeMemory:
      c.memSetStr(key, x)
      break
    case modeRedis:
      c.redisSetStr(key, x)
      break
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
}

func (c *cache) GetStr(key string) (string, bool) {
  switch(c.mode) {
    case modeMemory:
      return c.memGetStr(key)
    case modeRedis:
      return c.redisGetStr(key)
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
  return "", false
}

func (c *cache) Check(key string) bool {
  switch(c.mode) {
    case modeMemory:
      return c.memCheck(key)
    case modeRedis:
      return c.redisCheck(key)
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
  return false
}

func (c *cache) Remove(key string) {
  switch(c.mode) {
    case modeMemory:
      c.memRemove(key)
      break
    case modeRedis:
      c.redisRemove(key)
      break
    default:
      glog.Errorf("ERR: Mode Cache is Undefined\n")
  }
}

////
// Init
////
func New(mode string, expiryTime int64, redisURL string, redisMaxConnections int) *Cache {
  glog.Infof("LOG: CACHE: Init\n")
  glog.Infof("LOG: CACHE: Mode is %s\n", mode)
  glog.Infof("LOG: CACHE: Expiry Time = %d\n", expiryTime)
  if expiryTime == 0 {
		expiryTime = -1
	}
	c := &cache{
    items:             make(map[string]Item),
    itemsStr:          make(map[string]ItemStr),
		defaultExpiration: expiryTime,
	}
  switch mode {
    case "redis":
        conn, err := redis.DialURL(redisURL)
        if err == nil {
          conn.Close()
          c.redisPool = newRedisPool(redisURL, redisMaxConnections)
          c.mode = modeRedis
        } else {
          glog.Errorf("ERR: CACHE: REDIS: %v\n", err)
        }
        break
    case "memory":
        c.mode = modeMemory
        go c.memInit()
        break
  }
  return &Cache{c}
}

func (c *Cache) Close() {
  if c.redisPool != nil {
    c.redisPool.Close()
  }
}
