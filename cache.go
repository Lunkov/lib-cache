package cache

import (
  "time"
  "github.com/golang/glog"
)

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

////
// Init
////
func New(mode string, expiryTime int64, URL string, MaxConnections int) ICache {
  glog.Infof("LOG: CACHE: Init")
  glog.Infof("LOG: CACHE: Mode is %s", mode)
  glog.Infof("LOG: CACHE: Expiry Time = %d", expiryTime)
  if expiryTime == 0 {
		expiryTime = -1
	}
  switch mode {
    case "redis":
        return newRedis(mode, expiryTime, URL, MaxConnections)
    case "aerospike":
        return newAerospike(mode, expiryTime, URL, MaxConnections)
    case "map":
        return newMap(mode, expiryTime, URL, MaxConnections)
    case "syncmap":
        return newSyncMap(mode, expiryTime, URL, MaxConnections)
  }
  return nil
}


