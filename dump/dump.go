package dump

import (
	"encoding/json"

	"github.com/ricochhet/gomake/object"
)

func Dump(block object.FunctionBlock) (string, error) {
	if marshal, err := json.MarshalIndent(block, "", "\t"); err == nil {
		return string(marshal), nil
	} else {
		return "", err
	}
}
