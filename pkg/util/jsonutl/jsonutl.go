package jsonutl

import "encoding/json"

func ToJsonStr(o interface{}) string {
	vBytes, _ := json.Marshal(o)
	vJson := string(vBytes)
	return vJson
}
