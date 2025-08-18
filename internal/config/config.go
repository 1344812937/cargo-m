package config

import (
	"bytes"
	"github.com/creasty/defaults"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type ApplicationConfig struct {
	WebConfig       WebConfig       `toml:"http"`
	LocalRepoConfig LocalRepoConfig `toml:"maven_repo"`
}

type WebConfig struct {
	Host string `toml:"host" default:""`
	Port string `toml:"port" default:"8080"`
}

type LocalRepoConfig struct {
	LocalPath string `toml:"local_path" default:""`
}

func LoadApplicationConfig() *ApplicationConfig {
	configPath := "./cfg.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		panic("读取配置文件失败")
	}

	// 解析配置文件
	var cfg ApplicationConfig
	if err := toml.Unmarshal(content, &cfg); err != nil {
		panic("解析配置文件失败")
	}

	return &cfg
}

func createDefaultConfig(filePath string) *ApplicationConfig {
	defaultConfig := &ApplicationConfig{}
	err := defaults.Set(defaultConfig)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.SetIndentTables(true)
	encoder.SetTablesInline(false)
	err = encoder.Encode(defaultConfig)
	if err != nil {
		panic(err)
	}
	cfgToml := buf.String()
	println("配置：", cfgToml)
	create, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	_, writeErr := create.Write([]byte(cfgToml))
	if writeErr != nil {
		panic(writeErr)
	}
	return defaultConfig
}
