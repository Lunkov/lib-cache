package cache

import (
  "time"
  "sync"
  "encoding/json"
  "github.com/jasonlvhit/gocron"
  "github.com/golang/glog"
)

type CacheSyncMap struct {
  Cache
	items               sync.Map
}

func (c *CacheSyncMap) HasError() bool {
  return false
}

func (c *CacheSyncMap) GetMode() string {
  return "syncmap"
}

func (c *CacheSyncMap) Set(k string, x interface{}) {
  var e int64
	if c.defaultExpiration > 0 {
		e = time.Now().Add(time.Duration(c.defaultExpiration) * time.Second).UnixNano()
	} else {
    e = 0
  }
  c.items.Store(k, Item{
		Object:     x,
		Expiration: e,
	})
}

func (c *CacheSyncMap) Check(k string) bool {
  _, ok := c.items.Load(k)
  return ok
}

func (c *CacheSyncMap) Get(k string, obj interface{}) (interface{}, bool)  {
  i, ok := c.items.Load(k)
  if glog.V(9) {
    glog.Infof("DBG: CacheSyncMap: (ok = %v) item=%v", ok, i)
  }
  if !ok {
		return nil, false
	}
  if glog.V(9) {
    glog.Infof("DBG: CacheSyncMap: (ok = %v) item=%v", ok, i)
  }
  item := i.(Item)
  if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
  if glog.V(9) {
    glog.Infof("LOG: CacheSyncMap: Get(key = %s, ok = %v) item=%v", k, ok, item.Object)
  }
	return item.Object, true
}

func (c *CacheSyncMap) Remove(k string) {
  _, ok := c.items.Load(k)
  if ok {
    if glog.V(9) {
      glog.Infof("LOG: CacheSyncMap: Remove(key=%v) ok = %v", k, ok)
    }
    c.items.Delete(k);
  }
}

func (c *CacheSyncMap) Count() int64 {
  return 0
}

func (c *CacheSyncMap) Clear() {
  c.items = sync.Map{}
}

func (c *CacheSyncMap) ClearOld() {
  tn := time.Now().UnixNano()
  c.items.Range(func(key interface{}, value interface{}) bool {
    item := value.(Item)
    if item.Expiration != 0 && tn > item.Expiration {
      c.items.Delete(key)
    }
    return true
  })
}

func newSyncMap(mode string, expiryTime int64, uri string, maxConnections int) ICache {
 	c := &CacheSyncMap{
    Cache: Cache{ defaultExpiration:  expiryTime},
  }

  go c.initCron()
  
  return c
}

func (c *CacheSyncMap) initCron() {
  if c.defaultExpiration > 0 {
    s := gocron.NewScheduler()
    s.Every(1).Minutes().Do(c.ClearOld)
    <- s.Start()
  }
}

func (c *CacheSyncMap) GetAll2JSON(x interface{}) []byte {
  var memVals = make(map[string]interface{})
  res, _ := json.Marshal(memVals)
  return res
}
