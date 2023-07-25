package test

import (
	"fmt"
	"github.com/jason-wj/loadconfig/source/file"
	"testing"
)

func TestName(t *testing.T) {
	Setup(
		file.NewSource(file.WithPath("./settings.dev.yml")),
	)
	fmt.Println(AuthConfig.Secret)
}
