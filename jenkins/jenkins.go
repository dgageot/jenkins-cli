package jenkins

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(user, token, server, uri string) (io.ReadCloser, error) {
	var url = fmt.Sprintf("https://%s:%s@%s%s", user, token, server, uri)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	key, value, err := getCrumb(user, token, server)
	if err != nil {
		return nil, err
	}

	request.Header.Set(key, value)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func getCrumb(user string, token string, server string) (string, string, error) {
	var url = fmt.Sprintf("https://%v:%v@%v/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)", user, token, server)

	response, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return "", "", fmt.Errorf("http response code (%v)", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	crumbKeyValue := strings.Split(string(body), ":")
	return crumbKeyValue[0], crumbKeyValue[1], nil
}
