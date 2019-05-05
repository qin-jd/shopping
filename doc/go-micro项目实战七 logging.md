### logging
我们之前分别讲了微服务架构下的tracing,metrics。这节我们来讲一下logging

日志作为整个代码行为的记录，是程序执行逻辑和异常最直接的反馈。对于整个系统来说，日志是至关重要的组成部分。通过分析日志我们不仅可以发现系统的问题，同时日志中也蕴含了大量有价值可以被挖掘的信息，因此合理地记录日志是十分必要的。

个人理解 tracing对于系统性能分析很重要。
metrics对于系统的运行状态有了全局的把控。可以提前预警一些状态。
logging对于我们程序猿在遇到问题时具体去查错更加重要。

### logging日志库选型
go-micro默认的go-log，几乎不能看。pass了
去github上找到了两款最受欢迎的日志库组件。logrus（10.7K star），zap（6.6K star）
logrus:不支持日志切割，具有强大的hook功能。
zap:不支持日志切割，高性能，对。就是高性能。
对比之下，我选择了高性能。所以本节讲如何将zap集成到go-micro里

### zap使用指南
zap默认没有切割日志功能。官方推荐lumberjack来进行切割。
```
	fileName := "micro-srv.log"
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   128, //MB
		LocalTime: true,
		Compress:  true,
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(zap.DebugLevel))
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = log.Sugar()
```

### 总结
为了方便使用，我发布了个go-log包，把代码里的引用之前micro里go-log包替换成这个即可。只是没有logger.Log和Logf方法。有的都是zap里的方法例如Info和Infof,可以用此替代Log和Logf
https://github.com/qin-jd/go-log
