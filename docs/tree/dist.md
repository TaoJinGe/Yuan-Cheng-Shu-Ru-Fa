# dist TREE_MAP

server/ —— 服务端发布目录，保存 Linux 可执行文件和手机网页静态资源（生成产物，AI 不修改）
server/voice-bridge-server —— Linux amd64 服务端可执行文件，负责线上 HTTP 与 WebSocket 服务（生成产物，AI 不修改）
server/voice-bridge-server.exe —— Windows 本地服务端可执行文件，负责根目录启动脚本的本地 HTTP 与 WebSocket 服务（生成产物，AI 不修改）
server/static/ —— 手机网页发布资源目录，保存 index.html、app.js 和 style.css（生成产物，AI 不修改）
desktop/ —— Windows 客户端发布目录，保存桌面端 exe（生成产物，AI 不修改）
desktop/voice-bridge-client.exe —— Windows 客户端可执行文件，负责常驻连接与自动粘贴（生成产物，AI 不修改）
