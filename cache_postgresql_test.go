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

func TestCachePostgreSQL(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
    
  var info Person

  c1 := New("postgresql", 10, "postgre://dbuser:password@localhost:27017/testdb/keyvaluetest", 10)
  assert.Equal(t, "postgresql", c1.GetMode())
  c1.Clear()
  assert.Equal(t, false, c1.Check("1111"))

  uid, _ := uuid.Parse("00000002-0003-0004-1105-000000000001")
  info = Person{ID: uid, Login: "Max1", EMail: "max1@aaa.ru", Groups: []string{"g1", "g2"} }

  sessionToken := "f0g00002-0303-x004-1105-q00000000001"
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

  glog.Infof("LOG: postgresql.Get: (token = %v) user=%v", sessionToken, ipers15)
  glog.Infof("LOG: postgresql.Get: (token = %v) user=%v", sessionToken, ipers11)
  glog.Infof("LOG: postgresql.Get: (token = %v) user=%v", sessionToken, ipers)

  assert.Equal(t, true, err)
  
  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)
  
  assert.Equal(t, ipers.Groups, info.Groups)
  assert.Equal(t, int64(1), c1.Count())
  assert.Equal(t, true, c1.Check(sessionToken))

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


func BenchmarkPostgreSQL(b *testing.B) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "0")
	flag.Parse()
  
  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }

  c1 := New("postgresql", 10, "postgre://dbuser:password@localhost:27017/testdb/keyvaluetest", 10)
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
