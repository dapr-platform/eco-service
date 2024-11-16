package config

import "os"

var ECO_USER = ""
var ECO_APP_KEY = ""
var ECO_APP_SECRET = ""
var ECO_CALL_BACK_URL = ""
var ECO_OAUTH_CODE_URL = ""
var ECO_OAUTH_TOKEN_URL = ""
var ECO_OAUTH_CLIENT_ID = ""
var ECO_PROJECT_MAPPING = ""
var ECO_INVOKE_URL = ""
var ECO_PASSWORD = ""

var ECO_REALTIME_DATA_URL = ""

func init() {
	if v := os.Getenv("ECO_USER"); v != "" {
		ECO_USER = v
	}
	if v := os.Getenv("ECO_APP_KEY"); v != "" {
		ECO_APP_KEY = v
	}
	if v := os.Getenv("ECO_APP_SECRET"); v != "" {
		ECO_APP_SECRET = v
	}
	if v := os.Getenv("ECO_CALL_BACK_URL"); v != "" {
		ECO_CALL_BACK_URL = v
	}
	if v := os.Getenv("ECO_OAUTH_CODE_URL"); v != "" {
		ECO_OAUTH_CODE_URL = v
	}
	if v := os.Getenv("ECO_OAUTH_TOKEN_URL"); v != "" {
		ECO_OAUTH_TOKEN_URL = v
	}
	if v := os.Getenv("ECO_OAUTH_CLIENT_ID"); v != "" {
		ECO_OAUTH_CLIENT_ID = v
	}
	if v := os.Getenv("ECO_PROJECT_MAPPING"); v != "" {
		ECO_PROJECT_MAPPING = v
	}
	if v := os.Getenv("ECO_INVOKE_URL"); v != "" {
		ECO_INVOKE_URL = v
	}
	if v := os.Getenv("ECO_PASSWORD"); v != "" {
		ECO_PASSWORD = v
	}
	if v := os.Getenv("ECO_REALTIME_DATA_URL"); v != "" {
		ECO_REALTIME_DATA_URL = v
	}
}
