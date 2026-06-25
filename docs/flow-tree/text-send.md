# text-send FLOW_TREE

FLOW_ID：TEXT_SEND_001
功能名称：手机发送文本到同账号电脑当前光标处
入口页面：手机网页 /
触发点：登录后输入页的发送到电脑按钮
前端文件：/server/static/app.js、/desktop-client/internal/client/ws.go、/desktop-client/internal/paste/paste_windows.go
前端函数或事件：sendText、RunDesktopSocket、PasteText
API/后端入口：GET /ws 消息 send_text
后端核心文件：/server/internal/hub/hub.go
数据影响：records/账号/日期.md 文本记录文件，同一天发送内容追加到同一个 Markdown 文件
状态变化：手机端按下发送按钮时稳定触发一次发送，发送按钮贴近屏幕底部或输入法上沿，发送后保存上一条输入内容并自动清空输入框，顶部撤销按钮可恢复上一条继续编辑；按钮提供按下和发送闪烁反馈，并只把 message_type 为 send_text 的 ack 当作投递结果；服务端记录文本并向同账号所有桌面端发送 insert_text，记录文件按日期创建并按发送时间分段追加；桌面端有活动窗口时写剪贴板并 Ctrl+V，现代应用即使检测不到传统输入光标也会尝试粘贴；粘贴成功只更新客户端状态文字，不弹出系统通知
关联功能：AUTH_LOGIN_001、PRESENCE_SYNC_001
风险等级：高
验证方式：多台电脑同账号登录，当前输入框内有光标的电脑在手机发送后出现文本，手机输入框自动清空；点击顶部撤销后上一条文本回到输入框并可再次编辑发送；VS Code、Cursor、微信等现代输入框也应能粘贴
