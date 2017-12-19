package fastgo

import (
	"fmt"
	"io/ioutil"
	"errors"
	"strings"
)

var Conf = &Config{}

type Config struct {
	data map[string]string
}

func (c *Config) Prepare(configFile string) {
	c.data = make(map[string]string)
	err := c.Load(configFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("Config Prepare:", configFile)
}

func (c *Config) Load(configFile string) error {
	contentByte, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.New("cannot load config file")
	}
	content := string(contentByte)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "\r\t")
		if line == "" || line[0] == '#' {
			continue
		}
		splitData := strings.SplitN(line, "=", 2)
		if len(splitData) == 2 {
			for key, value := range splitData {
				splitData[key] = strings.Trim(value, "\r\t")
			}
			c.data[splitData[0]] = splitData[1]
		}
	}
	return nil
}

func (c *Config) Get(key string) string {
	if value, ok := c.data[key]; ok {
		return value
	} else {
		return ""
	}
}




