package config

import (
	"fmt"
	sdk "load-config"
	"load-config/source"
	"log"
)

var (
	_cfg *Settings
)

// Settings 兼容原先的配置结构
type Settings struct {
	Settings  Config `yaml:"settings"`
	callbacks []func()
}

func (e *Settings) runCallback() {
	for i := range e.callbacks {
		e.callbacks[i]()
	}
}

func (e *Settings) OnChange() {
	e.init()
	log.Println("!!! config change and reload")
}

func (e *Settings) Init() {
	e.init()
	log.Println("!!! config init")
}

func (e *Settings) init() {
	e.Settings.multiDatabase()
	e.runCallback()
}

// Config 配置集合
type Config struct {
	Application *Application          `yaml:"application"`
	Ssl         *Ssl                  `yaml:"ssl"`
	Auth        *Auth                 `yaml:"auth"`
	Database    *Database             `yaml:"database"`
	Databases   *map[string]*Database `yaml:"databases"`
	Gen         *Gen                  `yaml:"gen"`
}

// 多db改造
func (e *Config) multiDatabase() {
	if len(*e.Databases) == 0 {
		*e.Databases = map[string]*Database{
			"*": e.Database,
		}

	}
}

// Setup 载入配置文件
func Setup(s source.Source,
	fs ...func()) {
	_cfg = &Settings{
		Settings: Config{
			Application: ApplicationConfig,
			Ssl:         SslConfig,
			Auth:        AuthConfig,
			Database:    DatabaseConfig,
			Databases:   &DatabasesConfig,
			Gen:         GenConfig,
		},
		callbacks: fs,
	}
	var err error
	sdk.DefaultConfig, err = sdk.NewConfig(
		sdk.WithSource(s),
		sdk.WithEntity(_cfg),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("New config object fail: %s", err.Error()))
	}
	_cfg.Init()
}
