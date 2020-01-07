package network

type Client interface {
	NewRequest(url, method string, header map[string]string) Response
	PostFormData(body interface{}, url string) Response
	Post(body interface{}, url string) Response
	PostJSON(body interface{}, url string) Response
	Get(url string) Response
	GetWithoutJSON(url string) ([]byte, error)
}

type Response struct {
	Response map[string]interface{}
	Err      error
}
