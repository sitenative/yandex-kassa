package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// Client определяет транспортный уровень коммуникаций с API.
type Client struct {
	HTTP   *http.Client
	APIURL string
}

func (c *Client) request(ctx context.Context, method, urlStr string, body *[]byte) (*http.Response, error) {
	var request *http.Request

	u := fmt.Sprintf("%s/%s", c.APIURL, urlStr)

	var err error
	if body == nil {
		request, err = http.NewRequest(method, u, nil)
	} else {
		request, err = http.NewRequest(method, u, bytes.NewReader(*body))
	}
	if err != nil {
		return nil, err
	}

	return c.HTTP.Do(request.WithContext(ctx))
}
