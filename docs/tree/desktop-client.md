# desktop-client TREE_MAP

main.go —— Windows 客户端启动入口，只负责解析默认服务器参数并启动 GUI（修改影响客户端启动）
go.mod —— 客户端模块依赖文件，只负责声明构建依赖（修改影响依赖解析）
go.sum —— 客户端依赖校验文件，只负责锁定依赖校验值（修改影响依赖一致性）
app.manifest —— Windows 程序清单文件，只负责启用 Common Controls v6 和普通权限启动（修改影响 GUI 启动兼容性）
app.ico —— Windows 托盘图标文件，只负责窗口和右下角托盘显示（修改影响托盘可见性）
rsrc.syso —— Windows 资源文件，由 app.manifest 和 app.ico 生成并嵌入 exe（生成资源，AI 不手写修改）
internal/config/ —— 本地配置模块目录，只负责 config.json 的 token 读写（修改影响 30 天免登录）
internal/config/config.go —— 本地配置读写模块，只负责保存和清理 token 元数据及服务器地址（修改影响客户端自动登录，关联 FLOW_ID：AUTH_LOGIN_001、AUTH_LOGOUT_001）
internal/client/ —— 客户端联网模块目录，只负责登录 API 和 WebSocket 消息处理（修改影响连接状态与收文）
internal/client/api.go —— HTTP API 客户端，只负责注册、登录、退出和 token 校验请求（修改影响登录链路，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、AUTH_LOGOUT_001）
internal/client/ws.go —— WebSocket 客户端，只负责注册 desktop、接收 presence 和 insert_text（修改影响状态灯和插入链路，关联 FLOW_ID：PRESENCE_SYNC_001、TEXT_SEND_001）
internal/app/ —— 应用状态模块目录，只负责本地配置、认证状态和后台连接生命周期（修改影响 GUI 状态）
internal/app/app.go —— 应用状态协调模块，只负责注册登录、保存配置、启动/停止后台 WebSocket（修改影响常驻行为，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、PRESENCE_SYNC_001）
internal/ui/ —— Windows UI 模块目录，只负责登录注册窗口和托盘交互（修改影响桌面端界面）
internal/ui/icon.go —— Windows 图标加载模块，只负责加载内嵌图标并提供系统图标兜底（修改影响托盘可见性）
internal/ui/run.go —— 登录注册窗口模块，只负责窗口控件、按钮事件和状态展示（修改影响 GUI 登录体验，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001）
internal/ui/tray.go —— 系统托盘模块，只负责托盘菜单、显示窗口和退出软件（修改影响后台常驻体验，关联 FLOW_ID：DESKTOP_TRAY_001）
internal/ui/errorlog.go —— 启动错误提示模块，只负责 GUI 启动失败时写日志并弹窗提示（修改影响启动排错体验）
internal/paste/ —— Windows 粘贴模块目录，只负责剪贴板写入和 Ctrl+V 模拟（修改影响自动粘贴）
internal/paste/paste_windows.go —— Windows 自动粘贴实现，只负责检测输入光标、写剪贴板、SendInput 和恢复剪贴板（修改影响电脑光标处插入，关联 FLOW_ID：TEXT_SEND_001）
internal/paste/paste_other.go —— 非 Windows 占位实现，只负责阻止非 Windows 构建误用（修改影响跨平台构建提示）
