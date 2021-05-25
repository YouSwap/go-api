package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

func MD5(input string) (encode string, err error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return
	}
	encode = fmt.Sprintf("%X", md5.Sum(bytes))
	return
}
