#使用说明

## 切记，执行前需要确保自己的ip已经在公众号的白名单内。
## 白名单设置位于 
```go
微信公众后台 -> 基本配置（在左侧菜单底部） -> IP白名单
```
1. 获取公众号的appid和appsecert
2. 修改根目录下config.toml文件
    ```go
    appId = "wx5584eb1fcd64f74e"
    appSecret = "8515e6c8c73f1de7e5a3654cb52f86bf"
    ```
   将其修改为对应的值
3. 执行main.exe文件
    - 如果弹出的窗口打印“写入完成”，则表示获取成功
    - 如果没有任何打印，5秒后窗口关闭，则表示异常，需要将跟目录下wechat_openid.log文件打开查看。该文件记录异常信息
4. 用户的openid存储到openid.txt，自行查看