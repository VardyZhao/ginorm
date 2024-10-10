package config

import (
	"ginorm/cache"
	"ginorm/model"
	"ginorm/util"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Conf map[string]interface{} `yaml:"config"`
}

var Conf *Config

func LoadConfig(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var c Config
	c.Conf = make(map[string]interface{})
	err = yaml.Unmarshal(data, c.Conf)
	if err != nil {
		panic(err)
	}

	Conf = &c
}

func (c *Config) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	v := c.Conf

	for i := 0; i < len(keys); i++ {
		val, ok := v[keys[i]]
		if !ok {
			return nil
		}

		// 如果是最后一个键，直接返回值
		if len(keys)-1 == i {
			return val
		}

		// 如果不是最后一个键，检查是否为 map[string]interface{} 类型
		v, ok = val.(map[string]interface{})
		if !ok {
			return nil
		}
	}

	return nil
}

func (c *Config) GetString(key string) string {
	val := c.Get(key)
	var res string
	switch v := val.(type) {
	case int:
		res = strconv.Itoa(v) // 将整数转换为字符串
	case string:
		res = v
	default:
		return ""
	}
	return res
}

func (c *Config) GetInt(key string) int {
	val := c.Get(key)
	return val.(int)
}

// Init 初始化配置项
func Init() {
	// 从本地读取环境变量
	godotenv.Load()

	// 设置日志级别
	util.BuildLogger(os.Getenv("LOG_LEVEL"))

	// 读取翻译文件
	if err := LoadLocales("config/locales/zh-cn.yaml"); err != nil {
		util.Log().Panic("翻译文件加载失败", err)
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
	cache.Redis()
}
