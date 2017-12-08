# armyant(行军蚁)
[mqant](https://github.com/liangdas/mqant)压力测试工具
# 简介
armyant是从http压力测试工具[hey](https://github.com/rakyll/hey)改装而成。

hey只支持http接口的压力测试而armant可以自定义压测协议。

目前默认实现了http，mqtt两种协议的压力测试工具

# 依赖模块

    go get github.com/eclipse/paho.mqtt.golang

# 使用方法

> armyant无命令行工具,目前需要通过源码编译执行

### http压测

入口 http_task.go

具体实现: http_task/work.go

### mqtt压测

入口 mqtt_task.go

具体实现 mqtt_task/work.go


可以通过修改work.go代码来灵活更改具体的压测内容

### 操作系统最多文件打开限制

默认情况下普通操作系统都会限制系统同时打开的文件数量，mac系统默认是256.
如果不放开该限制armyant发出更多并发请求。

#### mac打开限制方式
  http://blog.csdn.net/mingtingjian/article/details/77675761

#### linux打开方式
  自己百度

#### mqant的压测参数

系统硬件:

   MAC电脑2核,16G内存,固态硬盘

进程:

 1. mqantserver 进程一个
 2. armyant 压测进程一个

压测结果:

    每一个连接每秒发出10个远程调用请求
    能达到的最大并发数为:3000

内存使用:

 mqantserver进程 90M

 CPU使用：

 mqantserver 120%

 aryant  100%

本次压测结果并不严谨,所用设备是自己的MAC电脑，同时还开启了很多编译器。
压测工具与测试进程也都在同一台机器,压测瓶颈主要在CPU性能上