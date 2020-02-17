package fasthttp_client

import (
	"encoding/json"
	"fmt"

	"jongme/app/network"

	"github.com/valyala/fasthttp"
)

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type FastHTTPClient struct {
	Client *fasthttp.Client
}

func (fhc *FastHTTPClient) NewRequest(url, method string, header map[string]string) network.Response {
	req := getJSONRequestWithHeader(url, method, header)
	fmt.Println("Request ==> ", req)
	resp := fasthttp.AcquireResponse()
	return do(req, resp, fhc.Client)
}

func (fhc *FastHTTPClient) PostFormData(body interface{}, url string) network.Response {
	req := getFormDataRequest(url, MethodPost)
	payloads, ok := body.(map[string]interface{})
	if !ok {
		return network.Response{Response: nil, Err: fmt.Errorf("Cannot cast payloads to []map[string]interface{}")}
	}

	var formBody string
	for key, value := range payloads {
		if formBody == "" {
			formBody = fmt.Sprintf("%s=%s", key, value)
		} else {
			formBody = fmt.Sprintf("%s&%s=%s", formBody, key, value)
		}
	}
	req.SetBodyString(formBody)

	resp := fasthttp.AcquireResponse()
	return do(req, resp, fhc.Client)
}

func (fhc *FastHTTPClient) Post(body interface{}, url string) network.Response {
	var bodyMessage []byte

	bodyMessage, _ = json.Marshal(body)

	req := getJSONRequest(url, MethodPost)
	req.SetBody(bodyMessage)
	fmt.Println("POST req: ", req)
	resp := fasthttp.AcquireResponse()

	return do(req, resp, fhc.Client)
}

func (fhc *FastHTTPClient) PostJSON(body interface{}, url string) network.Response {
	bodyMessage, _ := json.Marshal(body)

	req := getJSONRequest(url, MethodPost)
	req.SetBody(bodyMessage)
	resp := fasthttp.AcquireResponse()

	return do(req, resp, fhc.Client)
}

func (fhc *FastHTTPClient) Get(url string) network.Response {
	req := getJSONRequest(url, MethodGet)
	resp := fasthttp.AcquireResponse()
	return do(req, resp, fhc.Client)
}

func (fhc *FastHTTPClient) GetWithoutJSON(url string) ([]byte, error) {
	req := getJSONRequest(url, MethodGet)
	resp := fasthttp.AcquireResponse()
	return doWithoutJSON(req, resp, fhc.Client)
}

func doWithoutJSON(req *fasthttp.Request, resp *fasthttp.Response, c *fasthttp.Client) ([]byte, error) {
	err := c.Do(req, resp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if resp.StatusCode() != 200 {
		fmt.Println(string(resp.Body()))
		return nil, err
	}

	return resp.Body(), nil

}

func do(req *fasthttp.Request, resp *fasthttp.Response, c *fasthttp.Client) (r network.Response) {
	err := c.Do(req, resp)
	if err != nil {
		fmt.Println(err.Error())
		return network.Response{Response: nil, Err: err}
	}
	var response map[string]interface{}
	// response := map[string]interface{}{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
		return network.Response{Response: nil, Err: fmt.Errorf("%+v", err)}
	}
	if resp.StatusCode() != 200 {
		fmt.Println(string(resp.Body()))
		return network.Response{Response: response, Err: fmt.Errorf("Response status is %d", resp.StatusCode())}
	}

	// err = json.Unmarshal(resp.Body(), &response)
	// if err != nil {
	// 	fmt.Printf("ERROR: %+v\n", err)
	// 	return network.Response{Response: nil, Err: fmt.Errorf("%+v", err)}
	// }
	// fmt.Println("Response: %s\n", string(resp.Body()))

	return network.Response{Response: response, Err: nil}
}

func getJSONRequest(url string, method string) (req *fasthttp.Request) {
	req = fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json; charset=utf-8")
	return
}

func getFormDataRequest(url string, method string) (req *fasthttp.Request) {
	req = fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	return
}

func getJSONRequestWithHeader(url string, method string, header map[string]string) (req *fasthttp.Request) {
	req = fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")
	for key, value := range header {
		req.Header.Set(key, value)
	}
	return
}
