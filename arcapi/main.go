package arcapi

func GetPlayerUserMe(cookie string) (map[string]interface{}, error) {
	return requestAPI("user/me", cookie)
}
