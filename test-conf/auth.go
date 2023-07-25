package test

type Auth struct {
	Use     string
	Timeout int
	Secret  string
}

var AuthConfig = new(Auth)
