package cache

import (
  "time"
  "sync"
  "encoding/json"
  "github.com/jasonlvhit/gocron"
  "github.com/golang/glog"
)

type CacheMutexMap struct {
  Cache
	items               map[string]Item
  itemsStr            map[string]ItemStr
	mu                  sync.RWMutex
}

func (c *CacheMutexMap) HasError() bool {
  return false
}

func (c *CacheMutexMap) GetMode() string {
  return "mutexmap"
}

func (c *CacheMutexMap) Set(k string, x interface{}) {
  var e int64
	if c.defaultExpiration > 0 {
		e = time.Now().Add(time.Duration(c.defaultExpiration) * time.Second).UnixNano()
	} else {
    e = 0
  }
  c.mu.Lock()
  c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
  c.mu.Unlock()
}

func (c *CacheMutexMap) Check(k string) bool {
  c.mu.RLock()
  _, ok := c.items[k]
  c.mu.RUnlock()
  return ok
}

func (c *CacheMutexMap) Get(k string, obj interface{}) (interface{}, bool)  {
  c.mu.RLock()
  item, ok := c.items[k]
  if glog.V(9) {
    glog.Infof("DBG: memGet1: (ok = %v) item=%v\n", ok, item)
  }
  if !ok {
		c.mu.RUnlock()
		return nil, false
	}
  if glog.V(9) {
    glog.Infof("DBG: memGet2: (ok = %v) item=%v", ok, item)
  }
  if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
  if glog.V(9) {
    glog.Infof("LOG: memGet3: (ok = %v) item=%v", ok, item.Object)
    glog.Infof("LOG: memGet4: (ok = %v) item=%v", ok, obj)
  }
	return item.Object, true
}

func (c *CacheMutexMap) Remove(k string) {
  c.mu.Lock()
  _, ok := c.items[k]
  if ok {
    if glog.V(9) {
      glog.Infof("LOG: memRemove: (ok = %v) item=%v", ok, k)
    }
    delete(c.items, k);
  }
  c.mu.Unlock()
}

func (c *CacheMutexMap) Count() int64 {
  return int64(len(c.items))
}

func (c *CacheMutexMap) Clear() {
  for k := range c.items {
    delete(c.items, k)
  }
}

func (c *CacheMutexMap) ClearOld() {
  tn := time.Now().UnixNano()
  c.mu.Lock()
  for k, item := range c.items {
    if item.Expiration != 0 && tn > item.Expiration {
      delete(c.items, k)
    }
  }
  c.mu.Unlock()
}

func newMutexMap(mode string, expiryTime int64, uri string, maxConnections int) ICache {
 	c := &CacheMutexMap{
		Cache: Cache{ defaultExpiration:  expiryTime},
    items:                    make(map[string]Item),
    itemsStr:                 make(map[string]ItemStr),
	}

  go c.initCron()
  
  return c
}

func (c *CacheMutexMap) initCron() {
  if c.defaultExpiration > 0 {
    s := gocron.NewScheduler()
    s.Every(1).Minutes().Do(c.ClearOld)
    <- s.Start()
  }
}

func (c *CacheMutexMap) GetAll2JSON(x interface{}) []byte {
  c.mu.RLock()
  var memVals = make(map[string]interface{})
  for k, item := range c.items {
    memVals[k] = item.Object
  }
  c.mu.RUnlock()
  res, _ := json.Marshal(memVals)
  return res
}
