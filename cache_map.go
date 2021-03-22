package cache

import (
  "time"
  "encoding/json"
  "github.com/jasonlvhit/gocron"
  "github.com/golang/glog"
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type ItemStr struct {
	Object     string
	Expiration int64
}

type CacheMap struct {
  Cache
	items               map[string]Item
  itemsStr            map[string]ItemStr
}

func (c *CacheMap) HasError() bool {
  return false
}

func (c *CacheMap) GetMode() string {
  return "map"
}

func (c *CacheMap) Set(k string, x interface{}) {
  var e int64
	if c.defaultExpiration > 0 {
		e = time.Now().Add(time.Duration(c.defaultExpiration) * time.Second).UnixNano()
	} else {
    e = 0
  }
  c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

func (c *CacheMap) Check(k string) bool {
  _, ok := c.items[k]
  return ok
}

func (c *CacheMap) Get(k string, obj interface{}) (interface{}, bool)  {
  item, ok := c.items[k]
  if glog.V(9) {
    glog.Infof("DBG: memGet1: (ok = %v) item=%v\n", ok, item)
  }
  if !ok {
		return nil, false
	}
  if glog.V(9) {
    glog.Infof("DBG: memGet2: (ok = %v) item=%v\n", ok, item)
  }
  if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
  if glog.V(9) {
    glog.Infof("LOG: memGet3: (ok = %v) item=%v", ok, item.Object)
    glog.Infof("LOG: memGet4: (ok = %v) item=%v", ok, obj)
  }
	return item.Object, true
}

func (c *CacheMap) Remove(k string) {
  _, ok := c.items[k]
  if ok {
    if glog.V(9) {
      glog.Infof("LOG: memRemove: (ok = %v) item=%v", ok, k)
    }
    delete(c.items, k);
  }
}

func (c *CacheMap) Count() int64 {
  return int64(len(c.items))
}

func (c *CacheMap) Clear() {
  for k := range c.items {
    delete(c.items, k)
  }
}

func (c *CacheMap) ClearOld() {
  tn := time.Now().UnixNano()
  for k, item := range c.items {
    if item.Expiration != 0 && tn > item.Expiration {
      delete(c.items, k)
    }
  }
}

func newMap(mode string, expiryTime int64, uri string, maxConnections int) ICache {
 	c := &CacheMap{
		Cache: Cache{ defaultExpiration:  expiryTime},
    items:                    make(map[string]Item),
    itemsStr:                 make(map[string]ItemStr),
	}

  go c.initCron()
  
  return c
}

func (c *CacheMap) initCron() {
  if c.defaultExpiration > 0 {
    s := gocron.NewScheduler()
    s.Every(1).Minutes().Do(c.ClearOld)
    <- s.Start()
  }
}

func (c *CacheMap) GetAll2JSON(x interface{}) []byte {
  var memVals = make(map[string]interface{})
  for k, item := range c.items {
    memVals[k] = item.Object
  }
  res, _ := json.Marshal(memVals)
  return res
}
