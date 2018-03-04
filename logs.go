package k8s

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) GetLog(ctx context.Context, namespace string, res Resource, options ...Option) (string, error) {
	name := *res.GetMetadata().Name
	url, err := resourceGetURL(c.Endpoint, namespace, name, res, options...)
	if err != nil {
		return "", err
	}
	url += "/log"
	fmt.Println("getting", url)
	return c.getLog(ctx, url)
}

func (c *Client) getLog(ctx context.Context, url string) (log string, err error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %v", err)
	}
	if c.SetHeaders != nil {
		c.SetHeaders(r.Header)
	}

	re, err := c.client().Do(r)
	if err != nil {
		return "", fmt.Errorf("performing request: %v", err)
	}
	defer re.Body.Close()

	respBody, err := ioutil.ReadAll(re.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %v", err)
	}

	respCT := re.Header.Get("Content-Type")
	if err := checkStatusCode(respCT, re.StatusCode, respBody); err != nil {
		return "", err
	}
	log = string(respBody)
	return log, nil
}
