package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) NewRequest(method string, apiURL string, body io.Reader) (*http.Request, error) {
	apiPath, err := url.Parse(c.apiPath(apiURL))
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(apiPath)
	request, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Method": request.Method,
		"Accept": "application/json",
	}

	headers["Authorization"] = fmt.Sprintf("token %s", c.token)

	if body != nil {
		headers["Content-Type"] = "application/json"
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	for k, v := range c.headers {
		request.Header.Set(k, v)
	}

	values := request.URL.Query()
	request.URL.RawQuery = values.Encode()

	return request, nil
}

func (c *Client) RequestDecoder(method, path string, body io.Reader, v interface{}) error {
	request, err := c.NewRequest(method, path, body)
	if err != nil {
		return err
	}

	res, err := c.DoDecoder(request, v)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return err
}

func (c *Client) DoDecoder(req *http.Request, v interface{}) (*http.Response, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNoContent {
		return res, nil
	}

	err = checkErrorInResponse(res)
	if err != nil {
		fmt.Printf("Error sending request, %v", err)
		return res, err
	}

	if v != nil {
		var (
			resBuf bytes.Buffer
			resTee = io.TeeReader(res.Body, &resBuf)
		)
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resTee)
			return res, err
		}
		err = json.NewDecoder(resTee).Decode(v)
		if err != nil {
			fmt.Printf("Error parsing response %v", err)
		}
	}

	return res, err
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	response, err := c.c.Do(req)
	return response, err
}

func checkErrorInResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	return errors.New("Invalid Response")
}
