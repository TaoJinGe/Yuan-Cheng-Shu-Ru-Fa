# presence FLOW_TREE

FLOW_ID：PRESENCE_SYNC_001
功能名称：同步手机端和桌面端在线状态
入口页面：手机网页 /，Windows 客户端 GUI
触发点：WebSocket 连接、web_join、desktop_join、断开连接
前端文件：/server/static/app.js、/desktop-client/internal/ui/run.go、/desktop-client/internal/client/ws.go
前端函数或事件：connectSocket、resumeSocket、forceReconnect、startPing、afterLogin、reconnect、RunDesktopSocket、keepAlive
API/后端入口：GET /ws
后端核心文件：/server/internal/hub/hub.go
数据影响：服务端内存 web_connections 和 desktop_connections
状态变化：同账号手机端和至少一个桌面端在线时为已连接，否则为未连接或服务器断开；手机网页回到前台或点击顶部重连按钮时强制重建 WebSocket，并延迟 5 秒显示服务器断开；桌面端可点击重连按钮重建 WebSocket 且不需要重新输入密码；手机端发送业务 ping，桌面端发送 WebSocket 原生 ping 心跳减少空闲断线
关联功能：AUTH_LOGIN_001、TEXT_SEND_001
风险等级：高
验证方式：分别关闭手机网页和客户端，观察另一端红绿状态灯变化；手机浏览器切后台再切回前台应立即重建连接，点击顶部重连按钮不需要退出登录即可恢复连接
