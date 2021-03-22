package cache

import (
  "time"
  "fmt"
  "strings"
  "net"
  "net/url"
  
  "github.com/golang/glog"

  "encoding/json"
  
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
)

type ItemPostgreSQL struct {
  Key        string        `db:"id"                         json:"id"                                     gorm:"column:id;primary_key;"`
	UpdatedAt  time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: now()"     gorm:"type:timestamp with time zone"`
	Object     []byte        `db:"object"                     json:"object,ommitempty"                      `
}

type CachePostgreSQL struct {
  Cache
  MaxConnections int
  URL            string
  DBName         string
  TableName      string
  connect        *gorm.DB
}

func (c *CachePostgreSQL) HasError() bool {
  return c.connect == nil
}

func (c *CachePostgreSQL) GetMode() string {
  return "postgresql"
}

//
// PostgreSQL Cache
//

func (c *CachePostgreSQL) Set(key string, x interface{}) {
  if c.connect == nil {
    return
  }
	value, err := json.Marshal(x)
	if err != nil {
    glog.Errorf("ERR: POSTRESQL: JSON %s", err)
		return
	}
  var v0 ItemPostgreSQL
  v := ItemPostgreSQL{
    Key: key,
    Object: value,
  }
  tx := c.connect.Table(c.TableName).Begin()
  if err := tx.Where("id = ?", key).First(&v0).Error; err != nil {
    if gorm.IsRecordNotFoundError(err){
      tx.Create(&v)
    }
  } else {
    tx.Where("id = ?", key).Update(&v)
  }
  tx.Commit()
}

func (c *CachePostgreSQL) Check(key string) bool {
  if c.connect == nil {
    return false
  }
  var v ItemPostgreSQL
  c.connect.Table(c.TableName).Where("id = ?", key).First(&v)
  return v.Key == key
}

func (c *CachePostgreSQL) Get(key string, obj interface{}) (interface{}, bool) {
  if c.connect == nil {
    return "", false
  }
  var v ItemPostgreSQL
  c.connect.Table(c.TableName).Where("id = ?", key).First(&v)
  err := json.Unmarshal(v.Object, obj)
  if err != nil {
    glog.Errorf("ERR: CACHE: POSTRESQL: GET: %s", err)
    return nil, false
  }
  return obj, true
}

func (c *CachePostgreSQL) Remove(key string) {
  if c.connect == nil {
    return
  }
  v:= ItemPostgreSQL{Key: key}
  c.connect.Table(c.TableName).Delete(&v)
}

func (c *CachePostgreSQL) Clear() {
  if c.connect == nil {
    return
  }
  c.connect.Table(c.TableName).Where("1 = 1").Delete(&ItemPostgreSQL{})
}

func (c *CachePostgreSQL) Count() int64 {
  if c.connect == nil {
    return 0
  }
  count := int64(0)
  c.connect.Table(c.TableName).Select("count(id)").Count(&count)
  return count
}

func (c *CachePostgreSQL) GetAll2JSON(x interface{}) []byte {
  return []byte("[]")
}

func newPostgreSQL(mode string, expiryTime int64, uri string, maxConnections int) ICache {
  u, erru := url.Parse(uri)
  if erru != nil {
    glog.Errorf("ERR: CACHE: POSTGRESQL: URL(%s): %s", uri, erru)
    return nil
  }
  if u.Scheme != "postgre" {
    glog.Errorf("ERR: CACHE: POSTGRESQL: URL(%s): Scheme != aerospike", uri)
    return nil
  }
  namespace := u.Path
  if len(u.Path) > 1 {
    namespace = u.Path[1:]
  }
  ars := strings.Split(namespace, "/")
  if len(ars) != 2 {
    glog.Errorf("ERR: CACHE: POSTGRESQL: PARSE URL(%s) -> %v", uri, ars)
    return nil
  }
  host, port, _ := net.SplitHostPort(u.Host)
  pwd, _ := u.User.Password()
  
  connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
                          host, port, u.User.Username(), pwd, ars[0])

  db, err := gorm.Open("postgres", connectStr)
  if err != nil {
    glog.Errorf("ERR: CACHE: failed to connect database: %v", err)
    return nil
  }
  
  db.Table(ars[1]).AutoMigrate(ItemPostgreSQL{})
  
  return &CachePostgreSQL{
      connect: db,
      URL: connectStr,
      DBName: ars[0],
      TableName: ars[1],
    }
}

func (c *CachePostgreSQL) Close() {
  if c.connect != nil {
    c.connect.Close()
  }
}
