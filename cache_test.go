package cache

import (
  "testing"
  "github.com/stretchr/testify/assert"
  
  "bytes"
  "time"
  "encoding/gob"
  "encoding/json"

  "github.com/google/uuid"
  "github.com/golang/glog"
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
    glog.Errorf("ERR: Person: %s\n", err)
    return ""
  }
  return string(b)
}

func FromJSON(str string) (*Person, error) {
  p := Person{}
  if err := json.Unmarshal([]byte(str), &p); err != nil {
    glog.Errorf("ERR: Person: %s\n", err)
    return nil, err
  }
  return &p, nil
}

func TestCheckEnv(t *testing.T) {
  c1 := New("memory", 20, "", 10)
  assert.Equal(t, nil, c1)
}

func TestEMailUUID(t *testing.T) {
  uid, _ := uuid.Parse("abbf4958-17d9-56e3-afe4-30f21ebd1513")
  str := "login123123@mail.ru"
  id := uuid.NewSHA1(uuid.Nil, ([]byte)(str))
  assert.Equal(t, uid, id)
}

func TestLoadUnkonwn(t *testing.T) {

  c1 := New("unknown121212", 0, "", 10)
  assert.Equal(t, nil, c1)
}

func BenchmarkJSON(b *testing.B) {
  var p2 Person
  uid, _ := uuid.Parse("abbf4958-17d9-56e3-afe4-30f21ebd1513")
  p := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }
  for i := 0; i < b.N; i++ {
    buf, _ := json.Marshal(p)
    json.Unmarshal([]byte(buf), &p2)
    assert.Equal(b, p.Login, p2.Login)
  }
}

func BenchmarkGOB(b *testing.B) {
  var p2 Person
  uid, _ := uuid.Parse("abbf4958-17d9-56e3-afe4-30f21ebd1513")
  p := Person{ID: uid, Login: "Max", EMail: "max@aaa.ru", Groups: []string{"g1", "g2"} }
  buf := bytes.NewBuffer(nil)
  for i := 0; i < b.N; i++ {
    // encode
    enc := gob.NewEncoder(buf)
    enc.Encode(&p)
    // decode
    dec := gob.NewDecoder(buf)
    dec.Decode(&p2)
    assert.Equal(b, p.Login, p2.Login)
  }
}
