package config

import (
	"fmt"
	"load-config/source/file"
	"testing"
)

func TestName(t *testing.T) {
	Setup(
		file.NewSource(file.WithPath("./settings.dev.yml")),
	)
	fmt.Println(AuthConfig.Secret)
}
