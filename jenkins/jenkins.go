package jenkins

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func Get(user, token, server, query string) (body io.ReadCloser, err error) {
	key, value, err := crumb(user, token, server)
	if err != nil {
		return nil, err
	}

	uri := uri(user, token, server, query)
	request, err := http.NewRequest("GET", uri, nil)
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

func crumb(user, token, server string) (key string, value string, err error) {
	uri := uri(user, token, server, "/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)")

	response, err := http.Get(uri)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return "", "", errors.New(response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(string(body), ":")
	return parts[0], parts[1], nil
}

func uri(user, token, server, query string) (uri string) {
	return user + ":" + token + "@" + server + query
}
