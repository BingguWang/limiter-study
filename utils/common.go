package utils

import "encoding/json"

func ToJson(i interface{}) string {
	marshal, _ := json.Marshal(i)
	return string(marshal)
}
