package client

import (
	"crypto/md5"
	"eco-service/config"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type BoxResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Mac         string `json:"mac"`
		ProjectCode string `json:"projectCode"`
	} `json:"data"`
	Success bool `json:"success"`
}

func generateSign(data url.Values) string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += data.Get(k)
	}
	signStr += config.ECO_APP_SECRET

	fmt.Printf("Sign string before MD5: %s\n", signStr)
	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])
	fmt.Printf("Generated sign: %s\n", sign)
	return sign
}

func GetBoxProjectCode(mac string) (string, error) {
	fmt.Printf("Getting project code for MAC: %s\n", mac)
	if err := defaultClient.ensureValidToken(); err != nil {
		fmt.Printf("Failed to ensure valid token: %v\n", err)
		return "", err
	}

	// Prepare parameters
	data := url.Values{}
	data.Set("method", "GET_BOX")
	data.Set("client_id", config.ECO_APP_KEY)
	data.Set("access_token", defaultClient.accessToken)
	data.Set("timestamp", time.Now().Format("20060102150405"))
	data.Set("mac", mac)

	// Generate sign
	data.Set("sign", generateSign(data))

	fmt.Printf("Request parameters: %+v\n", data)

	// Create request
	req, err := http.NewRequest("POST", config.ECO_INVOKE_URL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := defaultClient.doRequest(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return "", err
	}

	fmt.Printf("Response body: %s\n", string(resp))

	var boxResp BoxResponse
	if err := json.Unmarshal(resp, &boxResp); err != nil {
		fmt.Printf("Failed to unmarshal response: %v\n", err)
		return "", err
	}
	if boxResp.Code != "0" {
		fmt.Printf("Box response: %+v\n", boxResp)
		return "", fmt.Errorf("invalid-mac")
	}
	fmt.Printf("Box response: %+v\n", boxResp)
	return boxResp.Data.ProjectCode, nil
}

func GetBoxesHourStats(params map[string]string) ([]byte, error) {
	fmt.Printf("Getting boxes hour stats with params: %+v\n", params)
	if err := defaultClient.ensureValidToken(); err != nil {
		fmt.Printf("Failed to ensure valid token: %v\n", err)
		return nil, err
	}

	// Prepare parameters
	data := url.Values{}
	data.Set("method", "GET_BOXES_HOUR_STATS")
	data.Set("client_id", config.ECO_APP_KEY)
	data.Set("access_token", defaultClient.accessToken)
	data.Set("timestamp", time.Now().Format("20060102150405"))

	// Add custom parameters
	for k, v := range params {
		data.Set(k, v)
	}

	// Generate sign
	data.Set("sign", generateSign(data))

	fmt.Printf("Request parameters: %+v\n", data)

	// Create request
	req, err := http.NewRequest("POST", config.ECO_INVOKE_URL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := defaultClient.doRequest(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return nil, err
	}

	fmt.Printf("Response body: %s\n", string(resp))
	return resp, nil
}
