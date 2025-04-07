package main

// import (
// 	"log/slog"

// 	"github.com/yiffyi/xfbbroker"
// 	"github.com/yiffyi/xfbbroker/xfb"
// )

// func main() {
// 	cfg = xfbbroker.LoadConfig()
// 	for k := range cfg.Users {
// 		u, ok := cfg.GetUser(k)
// 		if !ok {
// 			slog.Error("user not found", "name", k)
// 			continue
// 		}
// 		if u.Enabled {
// 			// get card info, balance
// 			s, newSessionId, err := xfb.GetUserDefaultLoginInfo(u.SessionId)
// 			if err != nil {
// 				slog.Error("unable to get user default login info", "err", err, "name", u.Name)
// 				continue
// 			}

// 			if newSessionId != "" {
// 				u.SessionId = newSessionId
// 			}

// 			balance, err := xfb.GetCardMoney(u.SessionId, u.YmUserId)
// 			if err != nil {
// 				slog.Error("unable to query card balance", "err", err, "name", u.Name)
// 				continue
// 			}
// 			if balance == "- - -" {
// 				slog.Info(`GetCardMoney returned "- - -"`)
// 				continue
// 			}

// 			slog.Info("Got user card info", "Username", u.Name, "Organization", s.SchoolName, "UserType", s.UserType, "Balance", balance)

// 		}
// 	}
// }
