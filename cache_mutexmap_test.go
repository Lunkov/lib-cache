package cache

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "flag"
  "fmt"
  "strconv"
  "github.com/google/uuid"
  "github.com/golang/glog"
  
  "github.com/Lunkov/lib-ref"
)

func TestLoadMutexMap(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

  var info Person

  c1 := New("mutexmap", 0, "", 10)
  c1.Clear()
  assert.Equal(t, int64(0), c1.Count())
  assert.Equal(t, false, c1.Check("1111"))

  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info = Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }

  sessionToken := "sdfsd-2345345-sdgfsdf--345"
  c1.Set(sessionToken, info)

  assert.Equal(t, int64(1), c1.Count())
  assert.Equal(t, true, c1.Check(sessionToken))

  var ipers Person
  ipers11, err := c1.Get(sessionToken, &info)
  
  if ref.GetType(ipers11) == "Person" {
    ipers, _ = ipers11.(Person)
  } else {
    // ipers, _ = *ipers11.(Person)
  }
  glog.Infof("LOG: c1.Get: (token = %v) user=%s\n", sessionToken, ref.GetType(ipers11))
  glog.Infof("LOG: c1.Get: (token = %v) user=%v\n", sessionToken, ipers11)
  glog.Infof("LOG: c1.Get: (token = %v) user=%v\n", sessionToken, ipers)

  assert.Equal(t, true, err)

  if &ipers == nil {
    t.Error(
      "For", "Person Convert",
      "expected", &info,
      "got", ipers,
    )
  }

  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)

  assert.Equal(t, info.Groups, ipers.Groups)

  // TODO
  //needJson := "{\"sdfsd-2345345-sdgfsdf--345\":{\"id\":\"00000002-0003-0004-0005-000000000001\",\"login\":\"Max\",\"email\":\"max@aaa.ru\",\"display_name\":\"\",\"avatar\":\"\",\"role\":\"\",\"groups\":[\"g1\",\"g2\"],\"TimeLogin\":\"0001-01-01T00:00:00Z\",\"AuthCode\":\"\"}}"  
  //resJson := string(c1.GetAll2JSON(Person{}))
  //assert.Equal(t, needJson, resJson)
  
  c1.Remove(sessionToken)

  assert.Equal(t, int64(0), c1.Count())
  assert.Equal(t, false, c1.Check(sessionToken))


  ipers1, err1 := c1.Get(sessionToken, info)

  assert.Equal(t, false, err1)
  assert.Equal(t, nil, ipers1)

}

func BenchmarkMutexMap(b *testing.B) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "0")
	flag.Parse()
  
  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }
  c1 := New("mutexmap", 0, "", 0)
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
