package cache

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "flag"
  "log"
  "fmt"
  "time"
  "github.com/google/uuid"
  "github.com/golang/glog"
  "encoding/json"
)

type Person struct {
  ID            uuid.UUID  `json:"id"            db:"id"`
  Login         string     `json:"login"         db:"login"`
  EMail         string     `json:"email"         db:"email"`
  DisplayName   string     `json:"display_name"  db:"display_name"`
  Avatar        string     `json:"avatar"        db:"avatar"`
  Role          string     `json:"role"          db:"srv_role"`
  Groups      []string     `json:"groups"        db:"groups"`
  TimeLogin     time.Time
  AuthCode      string
}

func (p *Person) ToJSON() string {
  b, err := json.Marshal(p)
  if err != nil {
    log.Printf("ERR: Person: %s\n", err)
    return ""
  }
  return string(b)
}

func FromJSON(str string) (*Person, error) {
  p := Person{}
  if err := json.Unmarshal([]byte(str), &p); err != nil {
    log.Printf("ERR: Person: %s\n", err)
    return nil, err
  }
  return &p, nil
}


func EqualStringArrays(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func TestCheckEnv(t *testing.T) {
  c1 := New("memory", 20, "", 10)
  res := c1.Mode()
  assert.Equal(t, "memory", res)
}

func TestEMailUUID(t *testing.T) {
  uid, _ := uuid.Parse("abbf4958-17d9-56e3-afe4-30f21ebd1513")
  str := "login123123@mail.ru"
  id := uuid.NewSHA1(uuid.Nil, ([]byte)(str))
  assert.Equal(t, uid, id)
}

func TestLoadUnkonwn(t *testing.T) {

  c1 := New("unknown121212", 0, "", 10)
  c1.Clear()
  assert.Equal(t, int64(-1), c1.Count())
  assert.Equal(t, "undefined", c1.Mode())
  assert.Equal(t, false, c1.Check("1111"))
  assert.Equal(t, true, c1.HasError())
}

func TestLoadMem(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  var info Person

  c1 := New("memory", 0, "", 10)
  c1.Clear()
  assert.Equal(t, int64(0), c1.Count())
  assert.Equal(t, false, c1.Check("1111"))

  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info = Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }

  var testCache = map[string]Person{}
  assert.Equal(t, 0, len(testCache))
  
  testCache["123"] = info
  assert.Equal(t, 1, len(testCache))

  sessionToken := "sdfsd-2345345-sdgfsdf--345"
  
  c1.memSet(sessionToken, info)

  assert.Equal(t, int64(1), c1.memCount())
  assert.Equal(t, int64(1), c1.Count())

  assert.Equal(t, true, c1.memCheck(sessionToken))

  ipers11, ok := c1.memGet(sessionToken, &info)
  ipers, _ := ipers11.(Person)


  assert.Equal(t, true, ok)

  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)

  assert.Equal(t, info.Groups, ipers.Groups)

  c1.memRemove(sessionToken)

  assert.Equal(t, int64(0), c1.memCount())
  assert.Equal(t, int64(0), c1.Count())

  assert.Equal(t, false, c1.memCheck(sessionToken))

  ipers1, ok1 := c1.memGet(sessionToken, &info)

  assert.Equal(t, false, ok1)
  assert.Equal(t, nil, ipers1)

  ipers2, ok2 := c1.Get(sessionToken, ipers1)

  assert.Equal(t, false, ok2)
  assert.Equal(t, nil, ipers2)
}

func TestLoadMemGlob(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()

  var info Person

  c1 := New("memory", 0, "", 10)
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
  
  if c1.GetType(ipers11) == "Person" {
    ipers, _ = ipers11.(Person)
  } else {
    // ipers, _ = *ipers11.(Person)
  }
  glog.Infof("LOG: c1.Get: (token = %v) user=%s\n", sessionToken, c1.GetType(ipers11))
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

  needJson := "{\"sdfsd-2345345-sdgfsdf--345\":{\"id\":\"00000002-0003-0004-0005-000000000001\",\"login\":\"Max\",\"email\":\"max@aaa.ru\",\"display_name\":\"\",\"avatar\":\"\",\"role\":\"\",\"groups\":[\"g1\",\"g2\"],\"TimeLogin\":\"0001-01-01T00:00:00Z\",\"AuthCode\":\"\"}}"  
  resJson := string(c1.GetAll2JSON(Person{}))
  
  assert.Equal(t, needJson, resJson)
  
  c1.Remove(sessionToken)

  assert.Equal(t, int64(0), c1.Count())
  assert.Equal(t, false, c1.Check(sessionToken))


  ipers1, err1 := c1.Get(sessionToken, info)

  assert.Equal(t, false, err1)
  assert.Equal(t, nil, ipers1)

}

func TestLoadRedis(t *testing.T) {
  var count_1 int64 = 1
  var info Person

  c1 := New("redis", 10, "redis://localhost:6379/0", 10)
  
  res := c1.Mode()
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
      "got", c1.redisCheck("1111"),
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
  if c1.GetType(ipers11) == "Person" {
    ipers, _ = ipers11.(Person)
  }
  if c1.GetType(ipers11) == "*Person" {
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

  if err != true {
    t.Error(
      "For", "Person UUID ERR",
      "expected", true,
      "got", err,
    )
  }
  
  assert.Equal(t, info.ID, ipers.ID)
  assert.Equal(t, info.Login, ipers.Login)
  assert.Equal(t, info.EMail, ipers.EMail)
  

  if !EqualStringArrays(ipers.Groups, info.Groups) {
    t.Error(
      "For", "Person Groups",
      "expected", info.Groups,
      "got", ipers.Groups,
    )
  }
  
  assert.Equal(t, count_1, c1.Count())
  assert.Equal(t, true, c1.Check(sessionToken))

  needJson := "{\"sdfsd-2345345-sdgfsdf--345\":{\"AuthCode\":\"\",\"TimeLogin\":\"0001-01-01T00:00:00Z\",\"avatar\":\"\",\"display_name\":\"\",\"email\":\"max1@aaa.ru\",\"groups\":[\"g1\",\"g2\"],\"id\":\"00000002-0003-0004-1105-000000000001\",\"login\":\"Max1\",\"role\":\"\"}}"  
  resJson := string(c1.GetAll2JSON(Person{}))
  if resJson != needJson {
    t.Error(
      "For", "GetAll2JSON: ",
      "expected", needJson,
      "got", resJson,
    )
  }

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
    //log.Printf("LOG: GET: %s => %v\n", code, ipers1)
    tp, _ := ipers1.(*Person)
    //log.Printf("LOG: GET: %s => %v\n", code, tp)
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

func BenchmarkMemory(b *testing.B) {
  uid, _ := uuid.Parse("00000002-0003-0004-0005-000000000001")
  info := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }
  c1 := New("memory", 0, "", 0)
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

    var ipers15 Person
    ipers11, _ := c1.Get(code, &ipers15)
    ipers, _ := ipers11.(*Person)
    if ipers != nil {
      b.Error(
        "For", "Person UUID",
        "expected", info,
        "got", ipers,
      )
    }

  }
 
  c1.Close()
}
