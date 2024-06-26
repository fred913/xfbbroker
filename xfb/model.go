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
