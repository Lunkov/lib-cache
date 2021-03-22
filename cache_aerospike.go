package cache

import (
  "strconv"
  "net"
  "net/url"
  "github.com/golang/glog"
  aero "github.com/aerospike/aerospike-client-go"
)

type CacheAerospike struct {
  Cache
  MaxConnections int
  URL            string
  Namespace      string
  wpolicy        *aero.WritePolicy
  rpolicy        *aero.BasePolicy
  connect        *aero.Client
}

func (c *CacheAerospike) HasError() bool {
  return c.connect == nil
}

func (c *CacheAerospike) GetMode() string {
  return "aerospike"
}

//
// Aerospike Cache
//

func (c *CacheAerospike) getKey(key string) *aero.Key {
  k, err := aero.NewKey(c.Namespace, c.GetMode(), key)
  if err != nil {
    glog.Errorf("ERR: CACHE: AEROSPIKE: getKey(%s): %v", key, err)
  }
  if glog.V(9) {
    glog.Infof("DBG: CACHE: AEROSPIKE: getKey(%s): %v", key, c.Namespace)
  }

  return k
}

func (c *CacheAerospike) Set(key string, x interface{}) {
  if c.connect == nil {
    return
  }
  err := c.connect.PutObject(c.wpolicy, c.getKey(key), x)
  if err != nil {
    glog.Errorf("ERR: AEROSPIKE: PUT(%s): %s", key, err)
    return
  }
}

func (c *CacheAerospike) Check(key string) bool {
  if c.connect == nil {
    return false
  }
  exists, err := c.connect.Exists(c.rpolicy, c.getKey(key))
  if err != nil {
    return false
  }
  return exists
}

func (c *CacheAerospike) Get(key string, obj interface{}) (interface{}, bool) {
  if c.connect == nil {
    return "", false
  }
  err := c.connect.GetObject(c.rpolicy, c.getKey(key), obj)
  if err != nil {
    glog.Errorf("ERR: CACHE: AEROSPIKE: GET(%s): %v", key, err)
    return nil, false
  }
  return obj, true
}

func (c *CacheAerospike) Remove(key string) {
  if c.connect == nil {
    return
  }
  c.connect.Delete(c.wpolicy, c.getKey(key))
}

func (c *CacheAerospike) Clear() {
  if c.connect == nil {
    return
  }
}

func (c *CacheAerospike) Count() int64 {
  if c.connect == nil {
    return 0
  }
  return 0
}

func (c *CacheAerospike) GetAll2JSON(x interface{}) []byte {
  return []byte("[]")
}

func newAerospike(mode string, expiryTime int64, uri string, maxConnections int) ICache {
  u, erru := url.Parse(uri)
  if erru != nil {
    glog.Errorf("ERR: CACHE: AEROSPIKE: URL(%s): %s", uri, erru)
    return nil
  }
  if u.Scheme != "aerospike" {
    glog.Errorf("ERR: CACHE: AEROSPIKE: URL(%s): Scheme != aerospike", uri)
    return nil
  }
  host, port, _ := net.SplitHostPort(u.Host)
  iport, _ := strconv.Atoi(port)
  
  policy := aero.NewClientPolicy()
  wpolicy := aero.NewWritePolicy(0, uint32(expiryTime))
  rpolicy := aero.NewPolicy()
  
  client, err := aero.NewClientWithPolicy(policy, host, iport)
  if err != nil {
    glog.Errorf("ERR: CACHE: AEROSPIKE: CONNECT: %s", err)
    return nil
  }
  namespace := u.Path
  if len(u.Path) > 1 {
    namespace = u.Path[1:]
  }
  return &CacheAerospike{
      connect: client,
      URL: uri,
      Namespace: namespace,
      wpolicy: wpolicy,
      rpolicy: rpolicy,
    }
}

func (c *CacheAerospike) Close() {
  if c.connect != nil {
    c.connect.Close()
  }
}
