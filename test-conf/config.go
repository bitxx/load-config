package test

import (
	"fmt"
	loadconfig "github.com/bitxx/load-config"
	"github.com/bitxx/load-config/source"
	"log"
)

// Config 配置集合
type Config struct {
	Application *Application          `yaml:"application"`
	Auth        *Auth                 `yaml:"auth"`
	Database    *Database             `yaml:"database"`
	Databases   *map[string]*Database `yaml:"databases"`
	Gen         *Gen                  `yaml:"gen"`
	callbacks   []func()
}

func (e *Config) runCallback() {
	for i := range e.callbacks {
		e.callbacks[i]()
	}
}

func (e *Config) OnChange() {
	e.init()
	log.Println("!!! config change and reload")
}

func (e *Config) Init() {
	e.init()
	log.Println("!!! config init")
}

func (e *Config) init() {
	e.multiDatabase()
	e.runCallback()
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
	_cfg := &Config{
		Application: ApplicationConfig,
		Auth:        AuthConfig,
		Database:    DatabaseConfig,
		Databases:   &DatabasesConfig,
		Gen:         GenConfig,
		callbacks:   fs,
	}
	var err error
	loadconfig.DefaultConfig, err = loadconfig.NewConfig(
		loadconfig.WithSource(s),
		loadconfig.WithEntity(_cfg),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("New config object fail: %s", err.Error()))
	}
	_cfg.Init()
}
