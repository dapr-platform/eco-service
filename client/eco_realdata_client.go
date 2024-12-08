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

type PowerMeterRealData struct {
	Elc2RepairType      int     `json:"elc2_repair_type"`
	TotalPower          float64 `json:"totalpower"`
	TotalPowerRepairType int    `json:"totalpower_repair_type"`
	Elc1ErrorType       int     `json:"elc1_error_type"`
	TotalRepairType     int     `json:"total_repair_type"`
	TotalRepair         float64 `json:"total_repair"`
	TotalPowerRepair    float64 `json:"totalpower_repair"`
	TotalPowerErrorType int     `json:"totalpower_error_type"`
	Total               float64 `json:"total"`
	Elc2Repair          float64 `json:"elc2_repair"`
	Elc1Repair          string  `json:"elc1_repair"`
	Elc2                float64 `json:"elc2"`
	Elc1                string  `json:"elc1"`
	TotalErrorType      int     `json:"total_error_type"`
	Elc1RepairType      int     `json:"elc1_repair_type"`
	Elc2ErrorType       int     `json:"elc2_error_type"`
	Timestamp           int64   `json:"timestamp"`
}


type PowerMeterData struct {
	Netport            int     `json:"netport"`
	Csq               int     `json:"csq"` 
	CsqErrorType      int     `json:"csq_error_type"`
	Syson             int     `json:"syson"`
	SbTimeRepairType  int     `json:"sb_time_repair_type"`
	EnergeRepairType  int     `json:"energe_repair_type"`
	NetportErrorType  int     `json:"netport_error_type"`
	SysonErrorType    int     `json:"syson_error_type"`
	Energe            float64 `json:"energe"`
	SbTimeRepair      int64   `json:"sb_time_repair"`
	EnergeErrorType   int     `json:"energe_error_type"`
	SbTime            int64   `json:"sb_time"`
	EnergeRepair      float64 `json:"energe_repair"`
	SysonRepairType   int     `json:"syson_repair_type"`
	SbTimeErrorType   int     `json:"sb_time_error_type"`
	SysonRepair       int     `json:"syson_repair"`
	CsqRepair         int     `json:"csq_repair"`
	NetportRepair     int     `json:"netport_repair"`
	NetportRepairType int     `json:"netport_repair_type"`
	CsqRepairType     int     `json:"csq_repair_type"`
	Timestamp         int64   `json:"timestamp"`
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
