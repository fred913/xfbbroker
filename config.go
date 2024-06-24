package xfbbroker

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
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
	lock                 sync.RWMutex `json:"-"`
	Users                []User
	LogFileName          string
	Debug                bool
	CheckTransInterval   int
	CheckBalanceInterval int
	ListenAddr           string
	AuthLocalUrl         string
	AuthCallback         string
}

func LoadConfig() *Config {
	content, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(content, &c)
	if err != nil {
		panic(err)
	}

	return &c
}

// need to be protected by lock
func (c *Config) save() {
	content, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("config.json", content, os.ModePerm)
	if err != nil {
		panic(err)
	}
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

func (c *Config) ReplaceUserBySessionId(session string, u *User) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k := range c.Users {
		if c.Users[k].SessionId == session {
			// u is already a copy of struct
			c.Users[k] = *u
			c.save()

			return nil
		}
	}

	return errors.New("could not found user with sessionId=" + session)
}

// func(*user) changed
func (c *Config) RWIterateUsers(fn func(*User) bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, u := range c.Users {
		if fn(&u) {
			// u is already a copy of struct
			c.Users[k] = u
		}
	}
	c.save()
}

func (c *Config) AppendUser(u *User) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Users = append(c.Users, *u)
	c.save()
}
