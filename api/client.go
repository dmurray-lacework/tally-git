package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	id      string
	token   string
	baseURL *url.URL
	c       *http.Client
	headers map[string]string

	Github *GithubService
}

type Option interface {
	apply(c *Client) error
}

type clientFunc func(c *Client) error

func (fn clientFunc) apply(c *Client) error {
	return fn(c)
}

func NewClient(apikey string) (*Client, error) {
	baseURL, err := url.Parse("https://api.github.com")
	if err != nil {
		return nil, err
	}

	c := &Client{
		id:      newID(),
		baseURL: baseURL,
		c:       &http.Client{},
		token:   apikey,
	}

	c.Github = &GithubService{c}

	return c, nil
}

func (c *Client) apiPath(path string) string {
	return fmt.Sprintf("%s", path)
}

func newID() string {
	now := time.Now().UTC().UnixNano()
	seed := rand.New(rand.NewSource(now))
	return strconv.FormatInt(seed.Int63(), 16)
}
