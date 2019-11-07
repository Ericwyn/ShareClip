# ShareClip
在多台设备/多个系统上面共享剪贴板

## 使用说明
  1. 首先你需要运行一个 ShareClipServer 作为一个 websocket 服务器
    - 使用 `./ShareClipServer -port 7878` 就可以启动了，默认监听的端口是7878
  2. 在你想要贡献剪贴板的设备上面运行 ShareClipClient
    - 使用 `./ShareClipClient -addr IP:PORT` 连接 ShareClipServer
  3. 之后所有连接 ShareClipServer 的设备都会同步更新剪贴板
  4. 支持 windows/linux/mac os
  

---

# ShareClip 
share your clipboard between windows/linux/mac

## Instruction
 1. Start a ShareClipServer on your pc
    - run the ShareClipServer
    - use `./ShareClipServer -port 7878` to set the listen port of ShareClipServer, the default port listen in 7878
 2. Start a ShareClipClient in any pc that your want to share clipboard with other devices
    - all the ShareClipClient should connect with one ShareClipServer
    - use `./ShareClipClient -addr IP:PORT` to connect the ShareClipServer
    - like `./ShareClipClient -addr 192.168.199.1:7878`
 3. any connect's devices will update the clipboard together
 4. support windows/linux/mac os

