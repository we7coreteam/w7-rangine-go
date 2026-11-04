# rangine-go

```
go get github.com/we7coreteam/w7-rangine-go/v2
```

## 日志上下文与轮转

`zap.Logger.With` 和 `slog.Logger.With` 的字段随日志转发到 stack 的各输出通道，嵌套 stack 与 slog 分组同样保留上下文。

`factory.Rotate("file")` 对已初始化的文件日志通道立即执行一次轮转：当前文件按时间戳备份，重新创建原路径文件继续写入，备份按 MaxBackups / MaxDays 清理。它不接收 stack 或 console 通道，也不会初始化尚未使用的通道。文件大小达到 MaxSize 时仍由 lumberjack 自动轮转。
