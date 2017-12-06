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

