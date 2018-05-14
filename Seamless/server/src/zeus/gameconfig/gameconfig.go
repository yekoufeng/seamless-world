package gameconfig

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

/*
 配置表工具
 获取完整配置表: config := gameconfig.New("test.json")
 获取单个值: v := config.Get("name")
 获取多值或者map: m := config.Get("game.match")
*/

// New 新建一个配置表
func New(file string) Config {
	return Config{file: file}
}

// Config 配置表
type Config struct {
	file string
	maps map[string]interface{}
}

// Get 获取配置
func (c *Config) Get(name string) interface{} {
	if c.maps == nil {
		c.read()
	}

	if c.maps == nil {
		return nil
	}

	keys := strings.Split(name, ".")
	l := len(keys)
	if l == 1 {
		return c.maps[name]
	}

	var ret interface{}
	for i := 0; i < l; i++ {
		if i == 0 {
			ret = c.maps[keys[i]]
			if ret == nil {
				return nil
			}
		} else {
			if m, ok := ret.(map[string]interface{}); ok {
				ret = m[keys[i]]
			} else {
				if l == i-1 {
					return ret
				}
				return nil
			}
		}
	}

	return ret
}

func (c *Config) read() {
	if !filepath.IsAbs(c.file) {
		file, err := filepath.Abs(c.file)
		if err != nil {
			panic(err)
		}
		c.file = file
	}

	bts, err := ioutil.ReadFile(c.file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bts, &c.maps)
	if err != nil {
		panic(err)
	}
}
