package cache

import (
  "reflect"
  "encoding/json"
  "github.com/golang/glog"
  "github.com/gomodule/redigo/redis"
)

type CacheRedis struct {
  Cache
  MaxConnections int
  URL            string
  pool           *redis.Pool
}

func (c *CacheRedis) HasError() bool {
  return c.pool == nil
}

func (c *CacheRedis) GetMode() string {
  return "redis"
}

//
// Redis Cache
//

func (c *CacheRedis) SetStr(key string, x string) {
  if c.pool == nil {
    return
  }
  var err error
  redisConn := c.pool.Get()
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

func (c *CacheRedis) Set(key string, x interface{}) {
  // serialize Object to JSON
	value, err := json.Marshal(x)
	if err != nil {
    glog.Errorf("ERR: REDIS: JSON %s\n", err)
		return
	}
  c.SetStr(key, string(value))
}

func (c *CacheRedis) Check(key string) bool {
  if c.pool == nil {
    return false
  }
  var ok int64 = 1
  redisConn := c.pool.Get()
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

func (c *CacheRedis) GetStr(key string) (string, bool) {
  if c.pool == nil {
    return "", false
  }
  redisConn := c.pool.Get()
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

func (c *CacheRedis) Get(key string, obj interface{}) (interface{}, bool) {
  data, ok := c.GetStr(key)
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

func (c *CacheRedis) Remove(key string) {
  if c.pool == nil {
    return
  }
  redisConn := c.pool.Get()
  redisConn.Do("DEL", key)
  redisConn.Flush()
  redisConn.Close()
}

func (c *CacheRedis) Clear() {
  if c.pool == nil {
    return
  }
  redisConn := c.pool.Get()
  keys, err := redis.Strings(redisConn.Do("KEYS", "*"))
  if err == nil {
    for _, key := range keys {
      redisConn.Do("DEL", key)
    }  
    redisConn.Flush()
  }
  redisConn.Close()
}

func (c *CacheRedis) Count() int64 {
  if c.pool == nil {
    return 0
  }
  var db_size int64 = 0
  redisConn := c.pool.Get()
  defer redisConn.Close()
  data, err := redisConn.Do("DBSIZE")
  if err != nil {
    glog.Infof("ERR: REDIS: Count: %v", err)
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

func (c *CacheRedis) GetAll2JSON(x interface{}) []byte {
  redisConn := c.pool.Get()
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

func newRedis(mode string, expiryTime int64, uri string, maxConnections int) ICache {
 	c := &CacheRedis{}
  conn, err := redis.DialURL(uri)
  if err == nil {
    conn.Close()
    c.pool = newRedisPool(uri, maxConnections)
  } else {
    glog.Errorf("ERR: CACHE: REDIS: %v", err)
  }
  return c
}

func newRedisPool(uri string, maxConnections int) *redis.Pool {
  if maxConnections < 1 {
    maxConnections = 100
  }
  glog.Infof("LOG: CACHE: REDIS (URL='%s', Max connections = %d)", uri, maxConnections)
  return &redis.Pool{
    // Maximum number of idle connections in the pool.
    MaxIdle: 80,
    // max number of connections
    MaxActive: maxConnections,
    // Dial is an application supplied function for creating and
    // configuring a connection.
    Dial: func() (redis.Conn, error) {
      c, err := redis.DialURL(uri)
      if err != nil {
        panic(err.Error())
      }
      return c, err
    },
  }
}

func (c *CacheRedis) Close() {
  if c.pool != nil {
    c.pool.Close()
  }
}
