package client

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"eco-service/config"
)

type TokenResponse struct {
	Data struct {
		TokenType    string `json:"tokenType"`
		AccessToken  string `json:"accessToken"`
		ExpiresIn    int64  `json:"expiresIn"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
	Success bool   `json:"success"`
	Code    string `json:"code"`
}

type EcoClient struct {
	client       *http.Client
	accessToken  string
	refreshToken string
	expireTime   time.Time
	mutex        sync.RWMutex
	stopRefresh  chan struct{}
}

var defaultClient *EcoClient

func init() {
	defaultClient = NewEcoClient()
}

func NewEcoClient() *EcoClient {
	client := &EcoClient{
		client:      &http.Client{},
		stopRefresh: make(chan struct{}),
	}
	go client.startRefreshTimer()
	return client
}

func (c *EcoClient) startRefreshTimer() {
	for {
		// Try to get initial token
		err := c.getInitialToken()
		if err != nil {
			fmt.Printf("Failed to get initial token: %v, retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Start refresh loop once we have initial token
		for {
			select {
			case <-c.stopRefresh:
				return
			case <-time.After(time.Until(c.expireTime) - 5*time.Minute):
				if err := c.refreshAccessToken(); err != nil {
					fmt.Printf("Failed to refresh token: %v, retrying in 5 seconds...\n", err)
					time.Sleep(5 * time.Second)
					// Try to get a new initial token if refresh fails
					if err := c.getInitialToken(); err != nil {
						fmt.Printf("Failed to get new token: %v, retrying in 5 seconds...\n", err)
						time.Sleep(5 * time.Second)
					}
				}
			}
		}
	}
}

func (c *EcoClient) getClientSecret(grantType, redirectUri, code string) string {
	str := config.ECO_APP_KEY + grantType + redirectUri + code + config.ECO_APP_SECRET
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

func (c *EcoClient) getInitialToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// First get auth code
	data := url.Values{}
	data.Set("response_type", "code")
	data.Set("client_id", config.ECO_APP_KEY)
	data.Set("redirect_uri", config.ECO_CALL_BACK_URL)
	data.Set("uname", config.ECO_USER)
	data.Set("passwd", config.ECO_PASSWORD)

	fmt.Printf("Getting auth code with data: %+v\n", data)
	fmt.Printf("Auth code URL: %s\n", config.ECO_OAUTH_CODE_URL)

	req, err := http.NewRequest("POST", config.ECO_OAUTH_CODE_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth code request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth code request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth code response body: %v", err)
	}
	fmt.Printf("Auth code response: %s\n", string(respBody))

	var codeResp struct {
		Code    string `json:"code"`
		Success string `json:"success"`
	}
	if err := json.Unmarshal(respBody, &codeResp); err != nil {
		return fmt.Errorf("failed to decode auth code response: %v, response body: %s", err, string(respBody))
	}

	if codeResp.Success != "true" {
		return fmt.Errorf("failed to get auth code, success=%s", codeResp.Success)
	}

	// Then exchange code for token
	tokenData := url.Values{}
	tokenData.Set("client_id", config.ECO_APP_KEY)
	tokenData.Set("grant_type", "authorization_code")
	tokenData.Set("redirect_uri", config.ECO_CALL_BACK_URL)
	tokenData.Set("code", codeResp.Code)
	tokenData.Set("client_secret", c.getClientSecret("authorization_code", config.ECO_CALL_BACK_URL, codeResp.Code))

	fmt.Printf("Getting token with data: %+v\n", tokenData)
	fmt.Printf("Token URL: %s\n", config.ECO_OAUTH_TOKEN_URL)

	tokenReq, err := http.NewRequest("POST", config.ECO_OAUTH_TOKEN_URL, strings.NewReader(tokenData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %v", err)
	}

	tokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenResp, err := c.client.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to send token request: %v", err)
	}
	defer tokenResp.Body.Close()

	tokenRespBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response body: %v", err)
	}
	fmt.Printf("Token response: %s\n", string(tokenRespBody))

	var tr TokenResponse
	if err := json.Unmarshal(tokenRespBody, &tr); err != nil {
		return fmt.Errorf("failed to decode token response: %v, response body: %s", err, string(tokenRespBody))
	}

	if !tr.Success {
		return fmt.Errorf("get token failed: %s", tr.Code)
	}

	c.accessToken = tr.Data.AccessToken
	c.refreshToken = tr.Data.RefreshToken
	c.expireTime = time.Now().Add(time.Duration(tr.Data.ExpiresIn) * time.Second)

	return nil
}

func (c *EcoClient) refreshAccessToken() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if time.Now().Before(c.expireTime) {
		return nil
	}

	data := url.Values{}
	data.Set("client_id", config.ECO_APP_KEY)
	data.Set("grant_type", "refresh_token")
	data.Set("redirect_uri", config.ECO_CALL_BACK_URL)
	data.Set("refresh_token", c.refreshToken)
	data.Set("client_secret", c.getClientSecret("refresh_token", config.ECO_CALL_BACK_URL, c.refreshToken))

	fmt.Printf("Refreshing token with data: %+v\n", data)

	req, err := http.NewRequest("POST", config.ECO_OAUTH_TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create refresh token request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send refresh token request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh token response body: %v", err)
	}
	fmt.Printf("Refresh token response: %s\n", string(respBody))

	var tokenResp TokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return fmt.Errorf("failed to decode refresh token response: %v, response body: %s", err, string(respBody))
	}

	if !tokenResp.Success {
		return fmt.Errorf("refresh token failed: %s", tokenResp.Code)
	}

	c.accessToken = tokenResp.Data.AccessToken
	c.refreshToken = tokenResp.Data.RefreshToken
	c.expireTime = time.Now().Add(time.Duration(tokenResp.Data.ExpiresIn) * time.Second)

	return nil
}

func (c *EcoClient) ensureValidToken() error {
	c.mutex.RLock()
	if c.accessToken != "" && time.Now().Before(c.expireTime) {
		c.mutex.RUnlock()
		return nil
	}
	c.mutex.RUnlock()

	if c.refreshToken == "" {
		return c.getInitialToken()
	}
	return c.refreshAccessToken()
}

func (c *EcoClient) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If unauthorized, try to refresh token and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		if err := c.getInitialToken(); err != nil {
			return nil, err
		}

		// Update Authorization header with new token
		req.Header.Set("Authorization", "Bearer "+c.accessToken)

		resp, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	return io.ReadAll(resp.Body)
}

func (c *EcoClient) Get(url string) ([]byte, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.accessToken)
	return c.doRequest(req)
}

func (c *EcoClient) Post(url string, body []byte) ([]byte, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.accessToken)
	req.Header.Add("Content-Type", "application/json")
	return c.doRequest(req)
}
