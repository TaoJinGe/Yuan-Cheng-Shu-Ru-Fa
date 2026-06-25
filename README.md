# 远程语音输入到电脑光标处

手机打开网页，用手机输入法语音转文字，点击“发送”后，同账号登录的 Windows 电脑会收到文本；只有当前处于输入状态、有输入光标的电脑会自动粘贴。

新版已支持账号注册/登录、30 天免登录、Windows 最小 GUI、右下角托盘常驻。

## 哪些放服务器，哪个在电脑打开

## 本地双击启动

如果只是在当前电脑本地测试，直接双击根目录：

```txt
启动本地服务.cmd
```

它会自动启动本地服务，并用默认浏览器打开：

```txt
http://127.0.0.1:8080/
```

本地账号数据库会保存在：

```txt
data/app.db
```

每次手机发送的文本会按账号保存到：

```txt
data/records/账号/日期.md
```

关闭本地服务，双击根目录：

```txt
关闭本地服务.cmd
```

它会自动关闭占用 `8080` 端口的本地服务。

如果网页里注册或登录一直显示“注册中 / 登录中”，通常表示本地服务没有运行。先双击 `启动本地服务.cmd`，确认浏览器地址是 `http://127.0.0.1:8080/`。

## 部署文件

上传到服务器：

```txt
dist/server/voice-bridge-server
dist/server/static/
```

如果你要把本地注册过的账号和发送记录一起迁移到宝塔，还需要上传：

```txt
data/app.db
data/records/
```

Windows 电脑双击打开：

```txt
dist/desktop/voice-bridge-client.exe
```

手机不用安装 App，浏览器打开：

```txt
https://srf.cccz.cc
```

## 宝塔部署

在宝塔文件管理里创建目录：

```txt
/www/wwwroot/srf.cccz.cc/
```

上传后目录应是：

```txt
/www/wwwroot/srf.cccz.cc/
  voice-bridge-server
  data/
    app.db              可选，本地账号数据库
    records/            可选，本地发送记录
  static/
    index.html
    app.js
    style.css
    favicon.ico
```

宝塔终端执行：

```bash
cd /www/wwwroot/srf.cccz.cc
chmod +x voice-bridge-server
./voice-bridge-server --port 8080 --db-file /www/wwwroot/srf.cccz.cc/data/app.db --records-dir /www/wwwroot/srf.cccz.cc/data/records
```

新开一个终端测试：

```bash
curl http://127.0.0.1:8080/health
```

返回 `ok` 就正常。

## 宝塔常驻运行

安装「Supervisor 管理器」，添加守护进程：

```txt
名称：voice-bridge
运行目录：/www/wwwroot/srf.cccz.cc
启动命令：/www/wwwroot/srf.cccz.cc/voice-bridge-server --port 8080 --db-file /www/wwwroot/srf.cccz.cc/data/app.db --records-dir /www/wwwroot/srf.cccz.cc/data/records
运行用户：root
```

账号不用提前创建。服务启动后，在 Windows 客户端或手机网页点击“注册”，输入账号密码即可。注册数据会进入 SQLite 数据库：

```txt
/www/wwwroot/srf.cccz.cc/data/app.db
```

手机每次发送的文本会保存到：

```txt
/www/wwwroot/srf.cccz.cc/data/records/账号/日期.md
```

同一个账号每天只生成一个 Markdown 文件，例如 `2026-06-25.md`。当天多次发送的文字会追加到同一个文件里，并用发送时间分段。

## Nginx 反向代理

宝塔添加站点：

```txt
srf.cccz.cc
```

反向代理目标：

```txt
http://127.0.0.1:8080
```

必须开启 WebSocket。手动配置可参考：

```nginx
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

建议在宝塔 SSL 里申请 Let's Encrypt，正式使用：

```txt
https://srf.cccz.cc
```

## Windows 客户端使用

双击：

```txt
dist/desktop/voice-bridge-client.exe
```

它会打开一个小窗口，不再弹 cmd。填写：

```txt
服务器 WebSocket 地址：wss://srf.cccz.cc/ws
账号：你的账号
密码：你的密码
[x] 30 天免登录
```

第一次使用点“注册”，以后点“登录”。登录成功后窗口会继续显示，方便确认连接状态；需要后台运行时，手动点击窗口右上角关闭按钮隐藏到右下角托盘。

托盘操作：

- 右键托盘图标 → 显示窗口
- 右键托盘图标 → 退出软件
- 点窗口右上角关闭按钮不会退出，只会隐藏到托盘
- 如果右下角没直接看到图标，点任务栏右下角的“隐藏的图标”小箭头查看

同一个账号可以在多台 Windows 电脑上同时登录。手机端发送后，所有同账号电脑都会收到消息；当前有活动窗口的电脑会尝试 Ctrl+V，所以发送前请先把目标电脑点进 VS Code、Cursor、记事本、微信输入框等输入区域。

## 手机使用

1. 手机浏览器打开 `https://srf.cccz.cc`。
2. 输入同一个账号和密码。
3. 没账号先点“注册”，已有账号点“登录”。
4. 点大文本框，用手机输入法语音转文字。
5. 先在目标电脑上点进 VS Code、Cursor、记事本、微信输入框等输入区域，让光标出现。
6. 手机点击“发送到电脑”。
7. 有输入光标的电脑会自动出现文字。

## 状态说明

- 未登录：没有有效 token。
- 红灯 / 未连接：只登录了一端，另一端不在线。
- 绿灯 / 已连接：手机网页和 Windows 客户端都在线。
- 服务器断开：客户端连不上服务器或 WebSocket 反代没配置好。

## 构建

服务端 Linux：

```powershell
cd server
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o ..\dist\server\voice-bridge-server .
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH
Copy-Item .\static\* ..\dist\server\static -Force
```

Windows GUI 客户端：

```powershell
cd desktop-client
go build -ldflags="-H windowsgui" -o ..\dist\desktop\voice-bridge-client.exe .
```

## 常见问题

### 双击没有 cmd 是正常的吗？

正常。新版是 GUI + 托盘程序。登录成功后窗口会继续显示；点击右上角关闭按钮后，它会在右下角后台运行。

### 账号密码保存在哪里？

服务端账号保存在 SQLite 数据库 `data/app.db`，密码是 bcrypt 哈希，不保存明文密码。客户端只保存 token，不保存明文密码。

### 30 天免登录为什么服务端重启后失效？

第一版 session 在内存中，服务端重启后旧 token 会失效，需要重新登录一次。账号不会丢，因为账号在 SQLite 的 `data/app.db`。

### 不能自动粘贴怎么办？

先把目标电脑光标点到 VS Code、Cursor、记事本、微信输入框等可输入区域。某些管理员权限窗口可能需要用相同权限启动客户端。新版会兼容一些检测不到传统输入光标的现代应用，但发送时仍建议只让目标电脑处于输入状态。
