package yaml

import (
	"github.com/bitxx/load-config/encoder"
	"github.com/bitxx/load-config/encoder/yaml/parser"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return parser.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return parser.Unmarshal(d, v)
}

func (y yamlEncoder) String() string {
	return "yaml"
}

func NewEncoder() encoder.Encoder {
	return yamlEncoder{}
}
