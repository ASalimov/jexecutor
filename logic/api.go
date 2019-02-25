package logic

import (
	"github.com/Jeffail/gabs"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"strings"
)

func req(method string, url string, body []byte, fn func(r *http.Request)) (int, *gabs.Container, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Add("Accept-Language", "en-us")
	fn(request)
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	//fmt.Println("salimov dd ",contents)
	if err != nil {
		return 0, nil, err
	}
	log.Debug("req: ", url, "\t data: ", string(body), "type:", response.Header.Get("Content-Type"), "\n rsp: ", string(contents))
	var rsp *gabs.Container
	if strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		rsp, err = gabs.ParseJSON(contents)
		if err != nil {
			return 0, nil, err
		}
	} else {
		rsp = nil
	}
	return response.StatusCode, rsp, nil
}
