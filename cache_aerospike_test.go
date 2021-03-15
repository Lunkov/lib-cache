package cache

import (
  "flag"
  "testing"
  "github.com/stretchr/testify/assert"
  
  "fmt"
  "strconv"
  "github.com/google/uuid"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-ref"
)


func TestCacheAerospike(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
    
  var info Person

  c1 := New("aerospike", 10, "aerospike://localhost:3001/test", 10)
  assert.Equal(t, "aerospike", c1.GetMode())
  c1.Clear()
  assert.Equal(t, false, c1.Check("1111"))

  uid, _ := uuid.Parse("00000002-0003-0004-1105-000000000001")
  info = Person{ID: uid, Login: "Max1", EMail: "max1@aaa.ru", Groups: []string{"g1", "g2"} }

  sessionToken := "sdfsd-2345345-sdgfsdf--345"
  c1.Set(sessionToken, info)
  //time.Sleep(2 * time.Second)

  var ipers15, ipers Person
  ipers11, err := c1.Get(sessionToken, &ipers15)
  // ipers, ok := ipers11.(*Person)
  if ref.GetType(ipers11) == "Person" {
    ipers, _ = ipers11.(Person)
  }
  if ref.GetType(ipers11) == "*Person" {
    ipers2, _ := ipers11.(*Person)
    ipers = *ipers2
  }

  glog.Infof("LOG: aerospike.Get: (token = %v) user=%v", sessionToken, ipers15)
  glog.Infof("LOG: aerospike.Get: (token = %v) user=%v", sessionToken, ipers11)
  glog.Infof("LOG: aerospike.Get: (token = %v) user=%v", sessionToken, ipers)

  assert.Equal(t, true, err)
  
  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)
  
  assert.Equal(t, ipers.Groups, info.Groups)
  assert.Equal(t, int64(0), c1.Count())
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
  
  assert.Equal(t, false, err3)
  assert.Equal(t, uuid.Nil, ipers4.ID)

  c1.Close()
}


func BenchmarkAerospike(b *testing.B) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "0")
	flag.Parse()
  
  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }

  c1 := New("aerospike", 0, "aerospike://localhost:3001/test", 100)
  c1.Clear()

  b.ResetTimer()
  for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
      b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
          code := fmt.Sprintf("code%d", i)
          c1.Set(code, info)
          var ipers5 Person
          _, ok := c1.Get(code, &ipers5)
          assert.Equal(b, true, ok)
        }
      })
    })
  }
  
  c1.Close()
  
}
