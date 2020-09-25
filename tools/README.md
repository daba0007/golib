# TOOLS
## 雪花算法

```Go
// datacenterid workerid
idGottor := NewSnowflake(0,0)
// 获取id
id := idGottor.NextVal()
// 获取数据中心ID和机器ID
datacenterid, workerid := GetDeviceID(id)
// 获取时间戳
t1 := GetTimestamp(id)
// 获取创建ID时的时间戳
t2 := GetGenTimestamp(id)
// 获取创建ID时的时间字符串(精度：秒)
s1 := GetGenTime(id)
// 获取时间戳已使用的占比：范围（0.0 - 1.0）
state := GetTimestampStatus()
```