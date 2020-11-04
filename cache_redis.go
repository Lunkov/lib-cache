package cache

import (
    "reflect"
    "encoding/json"
    "github.com/golang/glog"
    "github.com/gomodule/redigo/redis"
)

////
// Redis Cache
////
func (c *cache) redisSetStr(key string, x string) {
  var err error
  redisConn := c.redisPool.Get()
  if c.defaultExpiration > 0 {
    _, err = redisConn.Do("SETEX", key, int64(c.defaultExpiration), x)
    if glog.V(9) {
      glog.Errorf("DBG: REDIS: SETEX(%s, time=%d): item=%v : err%v\n", key, int64(c.defaultExpiration), x, err)
      return
    }
  } else {
    _, err = redisConn.Do("SET", key, x)
  }

  defer redisConn.Close()
  if err != nil {
    glog.Errorf("ERR: REDIS: SETEX(%s): %s\n", key, err)
    return
  }
  redisConn.Flush()
}

func (c *cache) redisSet(key string, x interface{}) {
  // serialize Object to JSON
	value, err := json.Marshal(x)
	if err != nil {
    glog.Errorf("ERR: REDIS: JSON %s\n", err)
		return
	}
  c.redisSetStr(key, string(value))
}

func (c *cache) redisCheck(key string) bool {
  var ok int64 = 1
  redisConn := c.redisPool.Get()
  d, err := redisConn.Do("EXISTS", key)
  defer redisConn.Close()
  if err != nil {
    return false
  }
  return d == ok
}

func fillStruct(data map[string]interface{}, result interface{}) {
    t := reflect.ValueOf(result).Elem()
    for k, v := range data {
        val := t.FieldByName(k)
        val.Set(reflect.ValueOf(v))
    }
}

func (c *cache) redisGetStr(key string) (string, bool) {
  redisConn := c.redisPool.Get()
  data, err := redis.String(redisConn.Do("GET", key))
  defer redisConn.Close()
  if err != nil {
    glog.Errorf("ERR: CACHE: REDIS: GET(%s): %s\n", key, err)
    return "", false
  }
  if glog.V(9) {
    glog.Infof("LOG: REDIS: GET: %s => %v\n", key, data)
  }
  return data, true
}

func (c *cache) redisGet(key string, obj interface{}) (interface{}, bool) {
  data, ok := c.redisGetStr(key)
  if !ok {
    glog.Errorf("ERR: CACHE: REDIS: GET(%s): !OK\n", key)
    return nil, false
  }
  err := json.Unmarshal([]byte(data), obj)
  if err != nil {
    glog.Errorf("ERR: CACHE: REDIS: GET: %s\n", err)
    return nil, false
  }
  return obj, true
}

func (c *cache) redisRemove(key string) {
  redisConn := c.redisPool.Get()
  redisConn.Do("DEL", key)
  redisConn.Flush()
  redisConn.Close()
}

func (c *cache) redisClear() {
  redisConn := c.redisPool.Get()
  keys, err := redis.Strings(redisConn.Do("KEYS", "*"))
  if err == nil {
    for _, key := range keys {
      redisConn.Do("DEL", key)
    }  
    redisConn.Flush()
  }
  redisConn.Close()
}

func (c *cache) redisCount() int64 {
  var db_size int64 = 0
  redisConn := c.redisPool.Get()
  defer redisConn.Close()
  data, err := redisConn.Do("DBSIZE")
  if err != nil {
    glog.Infof("ERR: Count: %v \n", err)
    return db_size
  }
  if data == nil {
    return db_size
  }
  if db_size, ok := data.(int64); ok {
    return db_size
  }
  return 0
}

func (c *cache) redisGetAll2JSON(x interface{}) []byte {
  redisConn := c.redisPool.Get()
  keys, err := redis.Strings(redisConn.Do("KEYS", "*"))
  defer redisConn.Close()
  if err != nil {
    glog.Errorf("ERR: redisGetLastValues: %s \n", err)
    return []byte("[]")
  }
  var memVals = make(map[string]interface{})
  var item interface{}
  for _, index := range keys {
    data, err := redis.String(redisConn.Do("GET", index))
    if err != nil {
    } else {
      err := json.Unmarshal([]byte(data), &item)
      if err != nil {
        glog.Errorf("ERR: CACHE: REDIS: GET: %s\n", err)
      } else {
        memVals[index] = item
      }
    }
  }
  res, _ := json.Marshal(memVals)
  return res
}

func newRedisPool(redisURL string, redisMaxConnections int) *redis.Pool {
  if redisMaxConnections < 1 {
    redisMaxConnections = 100
  }
  glog.Infof("LOG: CACHE: REDIS URL = %s\n", redisURL)
  glog.Infof("LOG: CACHE: REDIS Max connections = %d\n", redisMaxConnections)
  return &redis.Pool{
    // Maximum number of idle connections in the pool.
    MaxIdle: 80,
    // max number of connections
    MaxActive: redisMaxConnections,
    // Dial is an application supplied function for creating and
    // configuring a connection.
    Dial: func() (redis.Conn, error) {
      c, err := redis.DialURL(redisURL)
      if err != nil {
        panic(err.Error())
      }
      return c, err
    },
  }
}

