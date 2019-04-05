# 1 kate是什么

kate是一个轻量级HTTP API框架，可以让开发人员编写少量的代码来实现一个HTTP服务。

Kate也是一个library，提供和go标准HTTP完全兼容的router和middleware接口，可以嵌入标准HTTP或其他框架使用， 设计符合idiomatic Go。

此外，Kate包含了一些基础组件（redis，config，orm，eventbus等），封装了常用的lib， 以提升编程效率。

# 2 kate有哪些特点

## 2.1 支持context上下文

go语言的编程范式是生成大量的goroutine去执行具体的计算和IO任务，如果没有context的支持，很难维护goroutine间的派生关系。

比如我们处理一个请求，需要调用多个服务：passport、redis、mysql等，为了减少处理时间，我们对于可以并行的任务，都创建一个独立的goroutine去执行。

假如前端已经取消了请求，那与请求相关的所有子任务的结果实际上已经没有意义了，这时候如何取消这些任务，靠的就是context。

![](https://raw.githubusercontent.com/k81/kate/master/docs/resources/Understanding-Go-Context-Library-Google-Docs.png)

context维护了任务的派生关系树，当主任务取消或超时时，子任务会得到结束通知，这样实现了整个任务树的清理。

kate不只在请求的handler接口中支持context，日志模块也支持context。

下面的日志都属于同一个请求，处理请求的模块是httpsrv，请求ID是1。 （module=[httpsrv] session=[1]）。将来支持TraceID后，同一个请求的所有日志，都可以通过TraceID来关联。
	
     INFO    2017-04-22 20:15:50.850 [7527] msg=[request_in] module=[httpsrv] session=[1] remote=[172.16.3.201:43018] method=[POST] url=[/login/smsMt?source=yqh_merchant] body=[{ "username": "kaka", "ticket": "FL1TvN3OQL/+c+b6MZwkfeKpNi7tF4j3f7WXuCmWGDMcxMEKwyAMBuB3+c85mBAzk5cZIgOHONhmD6X03Qu9fAfeCDV5iIomTc6WM5sTNgRGHRWEhmDCvP0jsH/7c75+rdfPAmEhWF2KiadyXgAAAP//"}] fileline=[middleware_logging.go:14]
     DEBUG   2017-04-22 20:15:50.852 [7527] msg=[ticket parsed] module=[httpsrv] session=[1] username=[kaka] content=[{"i":4627242404091655169,"u":"kaka","c":1,"m":1,"s":"yqh_merchant","t":1492862908}] fileline=[ticket.go:77]
     DEBUG   2017-04-22 20:15:50.862 [7527] msg=[vcode key] module=[httpsrv] session=[1] username=[kaka] key=[kaka_24881055_1x9234#R234] fileline=[validation_code.go:65]
     DEBUG   2017-04-22 20:15:50.865 [7527] msg=[sending vcode sms] module=[httpsrv] session=[1] username=[kaka] cell=[18612965076] vcode=[2428] fileline=[sms.go:105]
     DEBUG   2017-04-22 20:15:51.114 [7527] msg=[sms sent] module=[httpsrv] session=[1] username=[kaka] cell=[18612965076] code=[2428] fileline=[sms.go:133]
     INFO    2017-04-22 20:15:51.123 [7527] msg=[request_out] module=[httpsrv] session=[1] status_code=[200] body=[{"status":10000,"msg":"success","data":[]}] duration_ms=[273] fileline=[middleware_logging.go:18]

## 2.2 支持middleware扩展

在kate中，请求的处理流程是这样的：
![](https://raw.githubusercontent.com/k81/kate/master/docs/resources/middlware.png)
中间件可以将通用的处理流程封装在一个地方，有效减少重复代码，提高编程效率。
## 2.3 配置集中化管理

传统的配置都保留在本地，这对于多节点部署的服务来说，更改一次配置，需要去多个机器上修改配置，很不方便。

kate的配置使用etcd服务，集中式管理。
![](https://raw.githubusercontent.com/k81/kate/master/docs/resources/etcd.png)

## 2.4 支持pprof性能监控

kate可以配置是否开启golang的http/pprof性能监控。pprof是个很好的工具，可以查看goroutine当前的执行栈，内存使用情况，GC等信息。

详见https://golang.org/pkg/net/http/pprof/
## 2.5 不停服重启

借助于gracehttp，可以实现不停服重启。

升级配置和程序，不会停服，从而大大服务可用性。
## 2.6 提供一些便利的组件，提高开发效率

kate提供如下组件：

    orm	        对db的面向对象封装。采用了和beego/orm完全一致的接口，去除了关联表和外键等影响性能的支持
    redismgr	redis manager, 链接池的封装，支持多codis代理节点负载均衡
    bigid	    64位ID生成器,支持分片因子继承
    singleflight  合并多个相同的并发查询，只发送一次真实资源请求，待取得结果后，返回所有调用方
    cache	    在redismgr基础上提供的缓存访问的简单接口
    retry	    重试策略库，提供多种重试策略：固定延迟重试、指数后退重试、限定最长时间的重试、限定最大次数的重试等，并支持策略组合
    redsync	    RedLock分布式锁的实现
    httpclient	http客户端组件
    taskengine	任务池组件，用于执行异步任务
    utils	    工具包，提供一些常用的工具函数
    config	    配置管理器，提供本地配置读取和etcd集中配置读取
    log	        日志组件

