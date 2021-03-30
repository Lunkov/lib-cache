package cache

import (
  "time"
  "github.com/golang/glog"
)

type CacheConfig struct {
  Mode            string     `yaml:"mode"`
  ExpiryTime      int64      `yaml:"expiry_time"`
  Url             string     `yaml:"url"`
  MaxConnections  int        `yaml:"max_connections"`
}

type ICache interface {
  HasError() bool
  GetMode() string
  
  Set(k string, obj interface{})
  Get(k string, obj interface{}) (interface{}, bool)
  Check(k string) bool
  Remove(k string)
  
  Clear()
  Count() int64
  
  Close()
}

const (
  // For use with functions that take an expiration time.
  NoExpiration time.Duration = -1
  // For use with functions that take an expiration time. Equivalent to
  // passing in the same expiration duration as was given to New() or
  // NewFrom() when the cache was created (e.g. 5 minutes.)
  DefaultExpiration time.Duration = 0
)

type Cache struct {
  defaultExpiration   int64
}

func (c *Cache) HasError() bool  { return true }
func (c *Cache) GetMode() string { return "undefined" }
func (c *Cache) Count() int64    { return 0 }
func (c *Cache) Set(k string, obj interface{}) {}
func (c *Cache) Get(k string, obj interface{}) (interface{}, bool) { return c, false }
func (c *Cache) Check(k string) bool {return false}
func (c *Cache) Clear() {}
func (c *Cache) Remove(k string) {}
func (c *Cache) Close() {}

//
// Init
//
func NewConfig(cfg *CacheConfig) ICache {
  return New(cfg.Mode, cfg.ExpiryTime, cfg.Url, cfg.MaxConnections)
}

func New(mode string, expiryTime int64, url string, maxConnections int) ICache {
  glog.Infof("LOG: CACHE: Init (mode=%s, ExpiryTime=%d)", mode, expiryTime)
  if expiryTime == 0 {
    expiryTime = -1
  }
  switch mode {
    case "map":
        return newMap(mode, expiryTime, url, maxConnections)
    case "syncmap":
        return newSyncMap(mode, expiryTime, url, maxConnections)
    case "redis":
        return newRedis(mode, expiryTime, url, maxConnections)
    case "aerospike":
        return newAerospike(mode, expiryTime, url, maxConnections)
// TODO
//    case "mongodb":
//        return newMongoDB(mode, expiryTime, url, maxConnections)
    case "postgresql":
        return newPostgreSQL(mode, expiryTime, url, maxConnections)

    case "mutexmap":
    default:
        return newMutexMap(mode, expiryTime, url, maxConnections)
  }
  return nil
}


