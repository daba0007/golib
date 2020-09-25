# CONFIG

config目前提供json文本的读取方式和yaml文本的读取方式

## json
```Go
// 定义自己结构体
type Config struct {
	DBCfg *DBConfig       `json:"DBService" yaml:"DBService"`
	ESCfg *[]ESConfig     `json:"EsService" yaml:"EsService"`
	GWCfg *Gateway        `json:"GatewayService" yaml:"GatewayService"`
	MQCfg *RabbitmqConfig `json:"RabbiqMQ" yaml:"RabbiqMQ"`
}
var cfg Config

// 对于json文本，调用ReadConfigFromJSONFile进行读取
// "config.json" -> 文本路径
ReadConfigFromJSONFile("config.json", &cfg)
// 对于json二进制流，调用ReadConfigFromJSONBytes进行读取
// data -> []byte
ReadConfigFromJSONBytes(data, &cfg)
```

## yaml
```Go
// 定义自己结构体
type Config struct {
	DBCfg *DBConfig       `json:"DBService" yaml:"DBService"`
	ESCfg *[]ESConfig     `json:"EsService" yaml:"EsService"`
	GWCfg *Gateway        `json:"GatewayService" yaml:"GatewayService"`
	MQCfg *RabbitmqConfig `json:"RabbiqMQ" yaml:"RabbiqMQ"`
}

var cfg Config
// 对于yaml文本，调用ReadConfigFromYAMLFile进行读取
// "config.yaml" -> 文本路径
ReadConfigFromYAMLFile("config.yaml", &cfg)
// 对于yaml二进制流，调用ReadConfigFromYAMLBytes进行读取
// data -> []byte
ReadConfigFromYAMLBytes(data, &cfg)
```