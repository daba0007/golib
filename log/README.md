# LOG

log包用于自定义日志库，来处理日志的输入输出

## 定义日志输出对象

```Go
// level -> log level DEBUG INFO WARNING ERROR FATAL
logger,err := NewLog(level, filepath)
```

## 日志输出

```Go
// level Debug
logger.Debug("Debug")
// level Info
logger.Info("Info")
// level Warning
logger.Warning("Warning")
// level Error
logger.Error("Error")
// level Fatal
logger.Fatal("Fatal")
```

## 设置日志库参数

```Go
// 修改日志保存时间,默认为7天
logger.SetSaveDays(30)
// 修改日志的输出等级
logger.SetLevel("INFO")
// 修改是否命令行输出，默认是
logger.SetConsoleFlag(false)
// 修改是否日志文件输出，默认是
logger.SetTextFlag(false)
```

## 关闭日志文件对象

```Go
logger.Close()
```