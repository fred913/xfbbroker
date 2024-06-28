package xfbbroker

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/yiffyi/xfbbroker/xfb"
)

type ApiServer struct {
	cfg *Config
}

func (s *ApiServer) probeSignPay(user *User) (string, error) {
	payUrl, err := xfb.RechargeOnCard("10.0", user.OpenId, user.SessionId, user.YmUserId)
	if err != nil {
		return "", err
	}

	u, _ := url.Parse(payUrl)
	tranNo := u.Query().Get("tran_no")
	_, err = xfb.SignPayCheck(tranNo)
	if err != nil {
		_, jumpUrl, err := xfb.GetSignUrl(tranNo)
		if err != nil {
			return "", err
		}

		return jumpUrl, nil
	}
	return "", nil
}

func (s *ApiServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	q := r.URL.Query()
	u, _ := url.Parse("https://auth.xiaofubao.com/auth/user/third/getCode")
	{
		q := u.Query()
		q.Set("callBackUrl", s.cfg.AuthCallback)
		u.RawQuery = q.Encode()
	}

	if q.Get("ymToken") == "" || q.Get("ymUserId") == "" {
		loc, err := xfb.GetRedirectLocation(u.String()) // Get the location: compatible with WeCom
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, loc, http.StatusTemporaryRedirect)
	} else {
		sess, data, err := xfb.GetUserById(q.Get("ymToken"), q.Get("ymUserId"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			u, ok := s.cfg.GetUser(data["id"].(string))
			if ok {
				u.SessionId = sess
				u.Failed = 0
				w.WriteHeader(http.StatusOK)
			} else {
				u = User{
					Name:      data["userName"].(string),
					OpenId:    data["thirdOpenid"].(string),
					SessionId: sess,
					YmUserId:  data["id"].(string),
					// Threshold: 100,
					Enabled: false,
				}
				w.WriteHeader(http.StatusCreated)
			}
			s.cfg.SetUser(u.YmUserId, u)
			s.cfg.Save()
		}
	}
}

func (s *ApiServer) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	q := r.URL.Query()
	sess := q.Get("sessionId")
	if len(sess) > 0 {
		user := s.cfg.SelectUserFromSessionId(sess)
		if user == nil {
			http.Error(w, "user with sessionId="+sess+" not found", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			body, err := json.MarshalIndent(user, "", "    ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}
	} else {
		http.Error(w, "no sessionId provided", http.StatusBadRequest)
	}
}

func (s *ApiServer) handleSignpay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	q := r.URL.Query()
	sess := q.Get("sessionId")
	if len(sess) > 0 {
		user := s.cfg.SelectUserFromSessionId(sess)
		if user == nil {
			http.Error(w, "user with sessionId="+sess+" not found", http.StatusNotFound)
			return
		}

		jumpUrl, err := s.probeSignPay(user)
		if err != nil {
			http.Error(w, "signPay check failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if len(jumpUrl) > 0 {
			w.Header().Set("Location", jumpUrl)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		http.Error(w, "no sessionId provided", http.StatusBadRequest)
	}
}

func CreateApiServer(cfg *Config) *mux.Router {
	r := mux.NewRouter()
	s := &ApiServer{
		cfg: cfg,
	}

	r.HandleFunc("/_/xfb/auth", s.handleAuth).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/_/xfb/signpay", s.handleSignpay).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/_/config", s.handleConfig).Methods(http.MethodGet, http.MethodPut, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	return r
}
