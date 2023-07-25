package test

import (
	"fmt"
"loadconfig/source/file"
"testing"
)

func TestName(t *testing.T) {
	Setup(
		file.NewSource(file.WithPath("./settings.dev.yml")),
	)
	fmt.Println(AuthConfig.Secret)
}
