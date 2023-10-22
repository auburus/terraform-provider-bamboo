package provider

import "net/http"

type BambooClient struct {
	BaseUrl  string
	Username string
	Password string
}

func NewBambooClient(baseUrl string, username string, password string) *BambooClient {
	return &BambooClient{baseUrl, username, password}
}

func (b *BambooClient) Request(method string, url string) (*http.Response, error) {
	// httpRequest, err := http.NewRequest(method, b.BaseUrl+url, nil)
	httpRequest, _ := http.NewRequest(method, b.BaseUrl+url, nil)
	httpRequest.SetBasicAuth(b.Username, b.Password)
	httpRequest.Header.Add("Accept", "application/json")

	return http.DefaultClient.Do(httpRequest)
}
