package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (client *TenableClient) Query(method, path, params string) (body []byte, err error) {
	url := fmt.Sprintf("https://cloud.tenable.com/%s", path)
	payload := strings.NewReader(params)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-ApiKeys", fmt.Sprintf("accessKey=%s;secretKey=%s", client.auth.AccessKey, client.auth.SecretKey))

	res, err := client.httpClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("query response code was: %d", res.StatusCode)
		return
	}

	body, _ = ioutil.ReadAll(res.Body)
	return
}
