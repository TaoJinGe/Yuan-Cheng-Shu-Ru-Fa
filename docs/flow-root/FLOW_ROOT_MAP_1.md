# FLOW_ROOT_MAP_1

认证会话链路 —— 手机网页登录注册、桌面客户端登录注册、token 校验、30 天免登录和退出登录。详见 /docs/flow-tree/auth.md（影响登录态、本地 token、WebSocket 鉴权）
设备在线状态链路 —— 手机网页端和桌面端 WebSocket 在线集合、presence 计算和红绿状态灯。详见 /docs/flow-tree/presence.md（影响用户判断两端是否连上）
远程文本发送链路 —— 手机输入文本、按 roomId 发送、服务端转发给桌面端、Windows 自动粘贴。详见 /docs/flow-tree/text-send.md（影响核心输入体验）
Windows 托盘常驻链路 —— 桌面端关闭窗口隐藏、托盘显示窗口和退出软件。详见 /docs/flow-tree/desktop-tray.md（影响后台常驻和退出体验）
