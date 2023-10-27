package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BambooClient struct {
	BaseUrl  string
	Username string
	Password string
}
type BambooProject struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewBambooClient(baseUrl string, username string, password string) *BambooClient {
	return &BambooClient{baseUrl, username, password}
}

func (b *BambooClient) Request(method string, url string, body io.Reader) (*http.Response, error) {
	httpRequest, _ := http.NewRequest(method, b.BaseUrl+url, body)
	httpRequest.SetBasicAuth(b.Username, b.Password)
	httpRequest.Header.Add("Accept", "application/json")
	if body != nil {
		httpRequest.Header.Add("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(httpRequest)
}

func (b *BambooClient) ReadProject(projectKey string) (*BambooProject, error) {
	httpResponse, err := b.Request("GET", "/project/"+projectKey, nil)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		defer httpResponse.Body.Close()

		body, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("Bamboo returned status code " + httpResponse.Status +
			"when looking for project " + projectKey +
			". Response:\n" + string(body))
	}

	var project BambooProject
	err = json.NewDecoder(httpResponse.Body).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (b *BambooClient) CreateProject(bambooProject *BambooProject) (*BambooProject, error) {
	var requestBody bytes.Buffer
	enc := json.NewEncoder(&requestBody)
	err := enc.Encode(bambooProject)

	if err != nil {
		return nil, err
	}

	httpResponse, err := b.Request("POST", "/project", &requestBody)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 201 {
		defer httpResponse.Body.Close()

		body, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("Bamboo returned status code " + httpResponse.Status +
			"when creating project " + fmt.Sprintf("%+v", bambooProject) +
			". Response:\n" + string(body))
	}

	var project BambooProject
	err = json.NewDecoder(httpResponse.Body).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (b *BambooClient) DeleteProject(projectKey string) error {
	httpResponse, err := b.Request("DELETE", "/project/"+projectKey, nil)
	if err != nil {
		return err
	}

	if httpResponse.StatusCode != 204 {
		defer httpResponse.Body.Close()

		body, err := io.ReadAll(httpResponse.Body)
		if err != nil {
			return err
		}

		return errors.New("Bamboo returned status code " + httpResponse.Status +
			"when deleting project " + projectKey +
			". Response:\n" + string(body))
	}

	return nil
}
