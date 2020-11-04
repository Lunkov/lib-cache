package cache

import (
  "time"
  "encoding/json"
  "github.com/jasonlvhit/gocron"
  "github.com/golang/glog"
)

func (c *cache) memSet(k string, x interface{}) {
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

func (c *cache) memCheck(k string) bool {
  c.mu.RLock()
  _, ok := c.items[k]
  c.mu.RUnlock()
  return ok
}

func (c *cache) memGet(k string, obj interface{}) (interface{}, bool)  {
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
    glog.Infof("DBG: memGet2: (ok = %v) item=%v\n", ok, item)
  }
  if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return nil, false
		}
	}
	c.mu.RUnlock()
  if glog.V(9) {
    glog.Infof("LOG: memGet3: (ok = %v) item=%v\n", ok, item.Object)
    glog.Infof("LOG: memGet4: (ok = %v) item=%v\n", ok, obj)
  }
	return item.Object, true
}

func (c *cache) memSetStr(k string, x string) {
  var e int64
	if c.defaultExpiration > 0 {
		e = time.Now().Add(time.Duration(c.defaultExpiration) * time.Second).UnixNano()
	} else {
    e = 0
  }
  c.mu.Lock()
  c.itemsStr[k] = ItemStr{
		Object:     x,
		Expiration: e,
	}
  c.mu.Unlock()
}

func (c *cache) memGetStr(k string) (string, bool)  {
  c.mu.RLock()
  item, ok := c.itemsStr[k]
  if !ok {
		c.mu.RUnlock()
		return "", false
	}
  if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.RUnlock()
			return "", false
		}
	}
	c.mu.RUnlock()
	return item.Object, true
}

func (c *cache) memRemove(k string) {
  c.mu.Lock()
  _, ok := c.items[k]
  if ok {
    if glog.V(9) {
      glog.Infof("LOG: memRemove: (ok = %v) item=%v\n", ok, k)
    }
    delete(c.items, k);
  }
  c.mu.Unlock()
}

func (c *cache) memCount() int64 {
  return int64(len(c.items))
}

func (c *cache) memClear() {
  for k := range c.items {
    delete(c.items, k)
  }
}

func (c *cache) memClearOld() {
  tn := time.Now().UnixNano()
  c.mu.Lock()
  for k, item := range c.items {
    if item.Expiration != 0 && tn > item.Expiration {
      delete(c.items, k)
    }
  }
  c.mu.Unlock()
}

func (c *cache) memInit() {
  if c.defaultExpiration > 0 {
    s := gocron.NewScheduler()
    s.Every(1).Minutes().Do(c.memClearOld)
    <- s.Start()
  }
}

func (c *cache) memGetAll2JSON(x interface{}) []byte {
  c.mu.RLock()
  var memVals = make(map[string]interface{})
  for k, item := range c.items {
    memVals[k] = item.Object
  }
  c.mu.RUnlock()
  res, _ := json.Marshal(memVals)
  return res
}
