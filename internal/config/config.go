package config

import (
	"bytes"
	"os"

	"github.com/creasty/defaults"
	"github.com/pelletier/go-toml/v2"
)

type ApplicationConfig struct {
	WebConfig       WebConfig       `toml:"http"`
	LocalRepoConfig LocalRepoConfig `toml:"maven_repo"`
	ProxyConfig     ProxyConfig     `toml:"proxy"`
}

type ProxyConfig struct {
	Enabled  bool   `toml:"enabled" default:"false"`
	Port     int    `toml:"port" default:"7890"`
	AuthUser string `toml:"auth_user" default:""`
	AuthPass string `toml:"auth_pass" default:""`
}

type WebConfig struct {
	Host string `toml:"host" default:""`
	Port string `toml:"port" default:"9080"`
}

type LocalRepoConfig struct {
	Enabled   bool   `toml:"enabled" default:"true"`
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
