package dap

import (
	"encoding/json"
	"errors"

	"github.com/raja-dettex/dap/schema"
)

type JSONEncoder struct {
}

type JSONDecoder struct{}

func (je *JSONEncoder) Encode(data any) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is nil")
	}
	return json.Marshal(data)
}

func (jd *JSONDecoder) Decode(raw []byte, data schema.Data) (schema.Data, error) {
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
