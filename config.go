package xfbbroker

import (
	"sync"

	"github.com/yiffyi/gorad/data"
)

type User struct {
	Name        string
	OpenId      string
	SessionId   string
	YmUserId    string
	Threshold   float64
	LastSerial  int
	WeComBotKey string
	Failed      int
	Enabled     bool
}

type Config struct {
	db                   *data.JSONDatabase
	lock                 *sync.RWMutex
	Users                map[string]User
	LogFileName          string
	Debug                bool
	CheckTransInterval   int
	CheckBalanceInterval int
	ListenAddr           string
	ListenTLS            bool
	TLSCertFile          string
	TLSKeyFile           string
	AuthLocalUrl         string
	AuthCallback         string
}

func LoadConfig() *Config {
	lock := sync.RWMutex{}
	db := data.NewJSONDatabase("config.json", true)
	cfg := Config{
		db:   db,
		lock: &lock,
	}

	db.Load(&cfg, true)
	return &cfg
}

func (c *Config) Save() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.db.Save(c)

}

func (c *Config) GetUser(k string) (User, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	u, ok := c.Users[k]
	return u, ok
}

func (c *Config) SetUser(k string, v User) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Users[k] = v
}

// return a copy of User
func (c *Config) SelectUserFromSessionId(session string) *User {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, u := range c.Users {
		if u.SessionId == session {
			// u is already a copy of struct
			return &u
		}
	}
	return nil
}
