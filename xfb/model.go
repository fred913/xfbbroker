package xfb

type XfbBaseResponse interface {
	GetStatusCode() int
}
type XfbResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	SessionId  string      `json:"-"`
}

func (r *XfbResponse) GetStatusCode() int {
	return r.StatusCode
}

type XfbQueryTransResponse struct {
	StatusCode int     `json:"statusCode"`
	Total      int     `json:"total"`
	Rows       []Trans `json:"rows"`
	Success    bool    `json:"success"`
}

func (r *XfbQueryTransResponse) GetStatusCode() int {
	return r.StatusCode
}

type Trans struct {
	Type           string `json:"type"`
	Time           string `json:"time"`
	Dealtime       string `json:"dealtime"`
	Address        string `json:"address"`
	FeeName        string `json:"feeName"`
	Serialno       string `json:"serialno"`
	Money          string `json:"money"`
	BusinessName   string `json:"businessName"`
	BusinessNum    string `json:"businessNum"`
	FeeNum         string `json:"feeNum"`
	AccName        string `json:"accName"`
	AccNum         string `json:"accNum"`
	PerCode        string `json:"perCode"`
	EWalletId      string `json:"eWalletId"`
	MonCard        string `json:"monCard"`
	AfterMon       string `json:"afterMon"`
	ConcessionsMon string `json:"concessionsMon"`
}

type UserDefaultLoginInfo struct {
	ID               string `json:"id"`
	SchoolCode       string `json:"schoolCode"`
	BadgeImg         string `json:"badgeImg"`
	SchoolName       string `json:"schoolName"`
	QrcodePayType    int    `json:"qrcodePayType"`
	UserName         string `json:"userName"`
	UserType         string `json:"userType"`
	JobNo            string `json:"jobNo"`
	UserIdcard       string `json:"userIdcard"`
	IdentityNo       string `json:"identityNo"`
	ThirdOpenid      string `json:"thirdOpenid"`
	ThirdbindMode    string `json:"thirdbindMode"`
	ThirdbindType    string `json:"thirdbindType"`
	CardPwdStatus    int    `json:"cardPwdStatus"`
	Sex              int    `json:"sex"`
	UserClass        string `json:"userClass"`
	RegiserTime      string `json:"regiserTime"`
	BindCardStatus   int    `json:"bindCardStatus"`
	BindCardTime     string `json:"bindCardTime"`
	HeadImg          string `json:"headImg"`
	DeviceId         string `json:"deviceId"`
	TestAccount      int    `json:"testAccount"`
	Token            string `json:"token"`
	Openid           string `json:"openid"`
	SchoolClasses    int    `json:"schoolClasses"`
	SchoolNature     int    `json:"schoolNature"`
	Platform         string `json:"platform"`
	CardPhone        string `json:"cardPhone"`
	BindCardRate     int    `json:"bindCardRate"`
	PayOpenid        string `json:"payOpenid"`
	CardIdentityType int    `json:"cardIdentityType"`
	AlumniFlag       int    `json:"alumniFlag"`
	PayType          int    `json:"payType"`
	CardActiveType   int    `json:"cardActiveType"`
	AuthType         int    `json:"authType"`
}
