package dump

import (
	"encoding/json"

	"github.com/ricochhet/gomake/object"
)

func Dump(block object.FunctionBlock) (string, error) {
	marshal, err := json.MarshalIndent(block, "", "\t")
	if err == nil {
		return string(marshal), nil
	}

	return "", err
}
