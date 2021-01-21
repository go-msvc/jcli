package httpcli

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/jcli"
)

func New(url string) (jcli.IClient, error) {
	return &httpClient{
		url: url,
	}, nil
}

type httpClient struct {
	url string
}

func (c *httpClient) Call(operName string, operReq jcli.IRequest, resType reflect.Type) (operRes jcli.IResponse, err error) {
	if val, ok := operReq.(jcli.IRequestValidator); ok {
		if err := val.Validate(); err != nil {
			return nil, errors.Wrapf(err, "invalid request")
		}
	}
	jsonOperReq, _ := json.Marshal(operReq)
	url := c.url + "/" + operName
	httpResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonOperReq))
	if err != nil {
		return nil, errors.Wrapf(err, "HTTP POST to %s failed", url)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		errorMessage, _ := ioutil.ReadAll(httpResp.Body)
		return nil, errors.Wrapf(err, "HTTP POST to %s -> %s %s", url, httpResp.Status, string(errorMessage))
	}

	jsonRes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "HTTP POST to %s failed to read response body", url)
	}
	newResValue := reflect.New(resType)
	if err := json.Unmarshal(jsonRes, newResValue.Interface()); err != nil {
		return nil, errors.Wrapf(err, "HTTP POST to %s failed to read response body as JSON", url)
	}
	return newResValue.Elem().Interface().(jcli.IResponse), nil
}
