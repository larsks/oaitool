package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type (
	ApiClient struct {
		AccessToken string
		ApiUrl      string
		client      *http.Client
	}
)

var ssourl = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"

func NewApiClient(offlinetoken string, apiurl string) *ApiClient {
	client := ApiClient{
		ApiUrl: apiurl,
		client: &http.Client{},
	}

	if err := client.GetAccessToken(offlinetoken); err != nil {
		return nil
	}

	return &client
}

func (client *ApiClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	log.Debugf("creating %s request for %s", method, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.AccessToken))

	return req, nil
}

func (client *ApiClient) GetAccessToken(offlinetoken string) error {
	var response TokenResponse

	params := url.Values{}
	params.Add("client_id", "cloud-services")
	params.Add("grant_type", "refresh_token")
	params.Add("refresh_token", strings.TrimSuffix(string(offlinetoken), "\n"))

	log.Debugf("asking %s for access token", ssourl)
	resp, err := client.client.PostForm(ssourl, params)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"failed to acquire token: %s",
			http.StatusText(resp.StatusCode),
		)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	client.AccessToken = response.AccessToken

	return nil
}
