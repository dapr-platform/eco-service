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

	"github.com/dapr-platform/common"
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

	hash := md5.Sum([]byte(signStr))
	sign := hex.EncodeToString(hash[:])
	return sign
}

func GetBoxProjectCode(mac string) (string, error) {
	common.Logger.Infof("Getting project code for MAC: %s\n", mac)
	if err := defaultClient.ensureValidToken(); err != nil {
		common.Logger.Errorf("Failed to ensure valid token: %v\n", err)
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

	// Create request
	req, err := http.NewRequest("POST", config.ECO_INVOKE_URL, strings.NewReader(data.Encode()))
	if err != nil {
		common.Logger.Errorf("Failed to create request: %v\n", err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := defaultClient.doRequest(req)
	if err != nil {
		common.Logger.Errorf("Request failed: %v\n", err)
		return "", err
	}

	var boxResp BoxResponse
	if err := json.Unmarshal(resp, &boxResp); err != nil {
		common.Logger.Errorf("Failed to unmarshal response: %v\n", err)
		return "", err
	}
	if boxResp.Code != "0" {
		respStr, _ := json.Marshal(boxResp)
		common.Logger.Errorf("Failed to get project code: %s\n", string(respStr))
		return "", fmt.Errorf("get-project-code-failed " + string(respStr))
	}
	return boxResp.Data.ProjectCode, nil
}
func GetBoxesMonthStats(params map[string]string) ([]byte, error) {
	return GetFunc("GET_BOXES_MON_STATS", params)

}
func GetFunc(method string, params map[string]string) ([]byte, error) {
	if err := defaultClient.ensureValidToken(); err != nil {
		common.Logger.Errorf("Failed to ensure valid token: %v\n", err)
		return nil, err
	}

	// Prepare parameters
	data := url.Values{}
	data.Set("method", method)
	data.Set("client_id", config.ECO_APP_KEY)
	data.Set("access_token", defaultClient.accessToken)
	data.Set("timestamp", time.Now().Format("20060102150405"))

	// Add custom parameters
	for k, v := range params {
		data.Set(k, v)
	}

	// Generate sign
	data.Set("sign", generateSign(data))

	// Create request
	req, err := http.NewRequest("POST", config.ECO_INVOKE_URL, strings.NewReader(data.Encode()))
	if err != nil {
		common.Logger.Errorf("Failed to create request: %v\n", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := defaultClient.doRequest(req)
	if err != nil {
		common.Logger.Errorf("Request failed: %v\n", err)
		return nil, err
	}

	return resp, nil
}

func GetBoxesHourStats(params map[string]string) ([]byte, error) {
	return GetFunc("GET_BOXES_HOUR_STATS", params)
}

func GetBoxesDayStats(params map[string]string) ([]byte, error) {
	return GetFunc("GET_BOXES_DAY_STATS", params)
}

func GetBoxesYearStats(params map[string]string) ([]byte, error) {
	return GetFunc("GET_BOXES_YEAR_STATS", params)
}
