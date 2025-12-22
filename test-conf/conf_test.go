package test

import (
	"fmt"
	loadconfig "github.com/bitxx/load-config"
	"github.com/bitxx/load-config/source/file"
	"testing"
)

func TestName(t *testing.T) {
	Setup(
		file.NewSource(file.WithPath("./settings.dev.yml")),
	)
	fmt.Println(loadconfig.DefaultConfig.Map())
}
