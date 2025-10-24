package config

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
)

func Load(filepath string, config any) error {
	file, err := os.Open(filepath);
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()
	return yaml.NewDecoder(file).Decode(config) //将YAML内容解码到配置中
}