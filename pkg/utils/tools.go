package utils

import (
	"encoding/json"
)

func JsonStr(data interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
