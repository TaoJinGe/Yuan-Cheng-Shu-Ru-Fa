# text-send FLOW_TREE

FLOW_ID：TEXT_SEND_001
功能名称：手机发送文本到同账号电脑当前光标处
入口页面：手机网页 /
触发点：登录后输入页的发送到电脑按钮
前端文件：/server/static/app.js、/desktop-client/internal/client/ws.go、/desktop-client/internal/paste/paste_windows.go
前端函数或事件：sendText、RunDesktopSocket、PasteText
API/后端入口：GET /ws 消息 send_text
后端核心文件：/server/internal/hub/hub.go
数据影响：records/账号/时间.md 文本记录文件
状态变化：服务端记录文本并向同账号所有桌面端发送 insert_text，只有检测到输入光标的桌面端写剪贴板并 Ctrl+V
关联功能：AUTH_LOGIN_001、PRESENCE_SYNC_001
风险等级：高
验证方式：多台电脑同账号登录，只有当前输入框内有光标的电脑在手机发送后出现文本
