package configfile

import (
	"encoding/json"
	"os"
)

func ReadConfig(path string, to any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, to)
}
