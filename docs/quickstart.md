本文将介绍如何快速搭建一个http服务。

## 1 准备工作

首选，需要配置好git ssh证书（请自行配置）

其次，配置GOPATH变量（参见01 Golang安装）
## 2 创建项目目录

Go项目的目录，一定要在GOPATH设定的目录中创建，目录层级如下：

    GOPATH
    └── src
        └── go.91power.com  //代码所属公司的域名
            └── tech-yqh     //项目组或部门名称
                └── example  //服务程序的名称（工程名，如passport-srv，order-srv等）

## 3 初始化项目代码

下载kate_init初始化工具，并在项目目录里执行。执行后，会自动在项目目录生成框架代码，

    #创建项目目录（正式项目，可以在git上创建好，再clone到本地）
    mkdir -p $GOPATH/src/go.example.com/tech/example
 
    #切换到项目目录（这个一定要做，kate_init脚本会根据当前执行目录的名字，来初始化程序的名字、配置等信息）
    cd $GOPATH/src/go.example.com/tech/example
 
    #初始化项目代码（首次执行会比较耗时间，请耐心等待，因为一些第三方package需要从github下载）
    curl -s  http://127.0.0.0.1/kate_init | bash -s

tree命令可以看下当前工程的目录结构
代码目录

    example
    ├── conf //配置目录
    │   ├── dev
    │   │   └── example.yaml  //开发环境配置文件
    │   ├── prod
    │   │   └── example.yaml  //生产环境配置文件
    │   └── qa
    │       └── example.yaml  //测试环境配置文件
    ├── httpsrv //http服务目录
    │   ├── config.go         //本模块配置项
    │   ├── errors.go         //错误码定义文件
    │   ├── handler_hello.go  //helloworld示例handler
    │   └── httpsrv.go        //模块入口文件，初始化router和http.Server对象
    ├── main.go   //程序入口文件
    ├── Makefile  //编译程序的Makefile，执行make编译
    ├── models    //数据模块，如程序不需访问DB，此模块可删除掉
    │   ├── config.go         //本模块配置项
    │   └── models.go         //模块入口文件，初始化db链接
    ├── service   //业务逻辑代码目录
    └── script    //脚本目录
        ├── build.sh          //编译和打包脚本
        ├── dev.sh            //开发环境配置导入脚本
        ├── prod.sh           //生产环境配置导入脚本
        └── qa.sh             //测试环境配置导入脚本


备注：现在错误码定义文件在httpsrv模块，这不是强制的，可以单独放到一个包里。
## 4 编译运行

    #编译项目，并使用dev配置
    #目前支持的环境变量有：dev,开发环境；qa，测试环境；preview, 预发环境;prod，生产环境。每个环境对应不同的数据库和etcd地址等配置信息。
    ./script/build.sh dev

执行后，会生成一个output目录，这个目录下的文件，是可以直接打包部署的。目录内容如下：

    output/
    ├── bin
    │   └── example //程序执行文件
    ├── conf
    │   ├── example.yaml //程序的本地配置文件
    │   └── gcached_example.yaml //从etcd配置服务加载的配置文件缓存，在etcd配置服务不可用时，程序仍能启动
    ├── log
    │   ├── example.err //错误日志文件
    │   └── example.log //日志文件
    └── run //pid文件目录。 程序启动后，会生成example.pid文件

查看程序版本信息
    
    #进入output目录
    cd output
    #查看程序版本号
    ./bin/example -version


接下来，导入配置，执行程序

    #导入开发配置，默认etcd服务器使用192.168.0.131节点
    #配置导入脚本有两种运行方式：不带参数，刷新所有配置；带参数，只刷新指定key的配置（如./dev.sh 'login/testusers'）
    ./script/dev.sh
 
    #进到打包输出目录output
    cd output
 
    #执行程序
    ./bin/example

## 6 浏览器测试接口

通过上面的步骤，程序默认会监听8000端口，可以用浏览器打开如下地址，验证下

    http://127.0.0.1:8000/hello

如下图：
![](https://raw.githubusercontent.com/k81/kate_docs/master/hello-world.png)
