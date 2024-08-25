package auth

type ApiRequest struct {
	baseUrl   string
	token     string
	source    string
	method    string
	timestamp int64
}

func NewApiRequest(baseUrl, token, source, method string, timestamp int64) *ApiRequest {
	return &ApiRequest{
		baseUrl:   baseUrl,
		token:     token,
		source:    source,
		timestamp: timestamp,
		method:    method,
	}
}
