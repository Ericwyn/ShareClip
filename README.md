# ShareClip
ShareClip 在多台设备/多个系统上面共享剪贴板

已在 windows10 和 ubuntu18.04 上测试通过

## 使用说明
ShareClip 通过一个服务端, 使得多个客户端可以进行剪贴板通信
### 服务端
使用以下命令可以启动一个服务端
```shell script
./ShareClip -server
```
server 支持的参数如下

 - `port` 监听的端口, 默认是 7878

### 客户端
使用以下命令可以启动一个服务端
```shell script
./ShareClip -client
```
server 支持的参数如下

 - `address` 监听的端口, 默认是 localhost:7878
 - `continue`  与服务断开连接之后是否持续重连, 默认只重连有限次数
 - `sender` 客户端标记名
 
### 共同支持的参数
 - `key`
 
    连接密码
 - `debug`
    
    开启 debug 日志
 - `v`
 
    查看版本号


## 注意事项

  1. ShareClip 只会共享你的文字剪贴, 无法共享文件/图片/视频的剪贴
