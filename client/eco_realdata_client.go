package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"eco-service/config"
)

type RealDataResponse[T any] struct {
	ResultCode string `json:"resultCode"`
	Message    Message  `json:"message"`
	QueryData  QueryData[T] `json:"queryData"`
}
type Message struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type QueryData[T any] struct {
	RtData T `json:"rtData"`
	State  any `json:"state"`
}
type WaterMeterData struct {
	MonthlyUsage          float64 `json:"monthlyuseage"`
	InsFlowRepairType     int     `json:"insflow_repair_type"`
	DailyUsageErrorType   int     `json:"dailyuseage_error_type"`
	BattVoltageRepair     float64 `json:"battVoltage_repair"`
	DailyUsage            float64 `json:"dailyuseage"`
	DailyUsageRepairType  int     `json:"dailyuseage_repair_type"`
	CumFlowErrorType      int     `json:"cumflow_error_type"`
	BattVoltageRepairType int     `json:"battVoltage_repair_type"`
	MonthlyUsageRepairType int    `json:"monthlyuseage_repair_type"`
	BattVoltage           float64 `json:"battVoltage"`
	InsFlowErrorType      int     `json:"insflow_error_type"`
	MonthlyUsageErrorType int     `json:"monthlyuseage_error_type"`
	DailyUsageRepair      float64 `json:"dailyuseage_repair"`
	InsFlowRepair         float64 `json:"insflow_repair"`
	CumFlow               float64 `json:"cumflow"`
	MonthlyUsageRepair    float64 `json:"monthlyuseage_repair"`
	CumFlowRepair         float64 `json:"cumflow_repair"`
	CumFlowRepairType     int     `json:"cumflow_repair_type"`
	BattVoltageErrorType  int     `json:"battVoltage_error_type"`
	InsFlow               float64 `json:"insflow"`
	Timestamp             int64   `json:"timestamp"`
}


func GetRealDataByCmCode[T any](cmCode string) (*RealDataResponse[T], error) {
	url := fmt.Sprintf("%s/dfm/device/realData?cmCode=%s", config.ECO_REALTIME_DATA_URL, cmCode)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get real data")
	}
	defer resp.Body.Close()

	var result RealDataResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "Failed to decode response")
	}

	return &result, nil
}
