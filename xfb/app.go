package xfb

func GetThirdUserAuthorize(ymToken, ymUserId string) {
	PostJSON(XfbApp+"/app/login/getThirdUserAuthorize", "", map[string]interface{}{
		"schoolCode": "20090820",
		"platform":   "WECHAT_H5",
		"ymToken":    ymToken,
		"ymUserId":   ymUserId,
	})
}
