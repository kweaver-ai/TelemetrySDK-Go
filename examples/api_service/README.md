# 典型API服务可观测埋点示例

original目录下为没有可观测数据埋点的原始代码，observable目录下为添加可观测数据埋点之后的代码。\
该示例有两个golang服务。portal_service为入口服务，接受外部请求，然后调用dependent_service的接口。

## trace埋点场景
1、gin server\
2、http client\
3、gorm client\
更多场景可参考 https://opentelemetry.io/ecosystem/registry/?s=&component=instrumentation&language=go

## 运行代码示例，上报可观测数据到AnyRobot
1、准备好已经安装AnyRobot Embedded 5的最新版本anyshare/anydata/anyfabric环境\
2、修改各个产品主模块中观测数据SDK的配置，开启观测数据记录
![telemetry-sdk-config](./images/telemetry-sdk-config.png)
3、在anyshare/anydata/anyfabric环境主机上面下载本示例代码，进入本文件所在目录执行下面命令下载依赖
```shell
go mod tidy
```
4、在anyshare/anydata/anyfabric环境主机上运行两个golang服务
```shell
go run observable/dependent_service/dependent_service_observability.go
go run observable/portal_service/portal_service_observability.go
```
5、模拟外部请求
```shell
curl localhost:50080/users/1
curl localhost:50080/users/2
```
6、登录AnyRobot Embedded 5仪表盘页面查看应用软件可观测仪表盘
