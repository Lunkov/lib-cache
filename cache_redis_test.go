package cache

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  
  "fmt"
  "github.com/google/uuid"
  "github.com/golang/glog"
)


func TestCacheRedis(t *testing.T) {
  var count_1 int64 = 1
  var info Person

  c1 := New("redis", 10, "redis://localhost:6379/0", 10)
  
  res := c1.GetMode()
  if res != "redis" {
    t.Error(
      "For", "Session Mode",
      "expected", "redis",
      "got", res,
    )
    return
  }
  c1.Clear()
  /*
  if c1.redisCount() != 0 {
    t.Error(
      "For", "Session Count Init",
      "expected", 0,
      "got", c1.redisCount(),
    )
  }*/
  if c1.Check("1111") != false {
    t.Error(
      "For", "Session Check Init",
      "expected", false,
      "got", c1.Check("1111"),
    )
  }

  uid, _ := uuid.Parse("00000002-0003-0004-1105-000000000001")
  info = Person{ID: uid, Login: "Max1", EMail: "max1@aaa.ru", Groups: []string{"g1", "g2"} }

  sessionToken := "sdfsd-2345345-sdgfsdf--345"
  c1.Set(sessionToken, info)
  //time.Sleep(2 * time.Second)

  var ipers15, ipers Person
  ipers11, err := c1.Get(sessionToken, &ipers15)
  // ipers, ok := ipers11.(*Person)
  if GetType(ipers11) == "Person" {
    ipers, _ = ipers11.(Person)
  }
  if GetType(ipers11) == "*Person" {
    ipers2, _ := ipers11.(*Person)
    ipers = *ipers2
  }

  glog.Infof("LOG: redis.Get: (token = %v) user=%v\n", sessionToken, ipers15)
  glog.Infof("LOG: redis.Get: (token = %v) user=%v\n", sessionToken, ipers11)
  glog.Infof("LOG: redis.Get: (token = %v) user=%v\n", sessionToken, ipers)

  /*
  if ok != true {
    t.Error(
      "For", "Person not OK",
      "expected", true,
      "got", ok,
    )
  }
*/
  // ipers, _ := FromJSON(*res2)

  assert.Equal(t, true, err)
  
  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)
  
  assert.Equal(t, ipers.Groups, info.Groups)
  assert.Equal(t, count_1, c1.Count())
  assert.Equal(t, true, c1.Check(sessionToken))

  //needJson := "{\"sdfsd-2345345-sdgfsdf--345\":{\"AuthCode\":\"\",\"TimeLogin\":\"0001-01-01T00:00:00Z\",\"avatar\":\"\",\"display_name\":\"\",\"email\":\"max1@aaa.ru\",\"groups\":[\"g1\",\"g2\"],\"id\":\"00000002-0003-0004-1105-000000000001\",\"login\":\"Max1\",\"role\":\"\"}}"  
  // TODO
  //resJson := string(c1.GetAll2JSON(Person{}))
  //assert.Equal(t, needJson, resJson)

  c1.Remove(sessionToken)
  //time.Sleep(2 * time.Second)
  assert.Equal(t, int64(0), c1.Count())
  assert.Equal(t, false, c1.Check(sessionToken))

  var ipers5 *Person
  ipers11, err3 := c1.Get(sessionToken, ipers5)
  ipers4, _ := ipers11.(Person)
  //ipers4, _ := FromJSON(*res3)
  
  
  if err3 != false {
    t.Error(
      "For", "Person UUID ERR",
      "expected", false,
      "got", err,
    )
  }
  
  assert.Equal(t, uuid.Nil, ipers4.ID)

  c1.Close()
}


func BenchmarkRedis(b *testing.B) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "0")
	flag.Parse()
  
  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }

  c1 := New("redis", 0, "redis://localhost:6379/0", 100)
  c1.Clear()

  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    count := c1.Count()
    count ++
  }
  
  for i := 0; i < b.N; i++ {
    code := fmt.Sprintf("code%d", i)
    c1.Set(code, info)
  }

  count := c1.Count()
  if count != int64(b.N) {
    b.Error(
      "For", "Count After Add",
      "expected", int64(b.N),
      "got", count,
    )
  }

  for i := 0; i < b.N; i++ {
    code := fmt.Sprintf("code%d", i)
    var ipers5 Person
    ipers1, ok := c1.Get(code, &ipers5)
    if !ok {
      b.Error(
        "For", "c1.Get",
        "expected", true,
        "got", ok,
      )
    }
    tp, _ := ipers1.(*Person)
    if tp.ID != info.ID {
      b.Error(
        "For", "Person UUID",
        "expected", info.ID,
        "got", tp.ID,
      )
    }
  }
 
  c1.Close()
  
}
