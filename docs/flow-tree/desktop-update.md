# desktop-update FLOW_TREE

FLOW_ID：DESKTOP_UPDATE_001
功能名称：Windows 客户端检查更新
入口页面：Windows 客户端登录窗口
触发点：客户端启动并创建主窗口后，或点击检查更新按钮
前端文件：/desktop-client/internal/ui/run.go、/desktop-client/internal/update/checker.go
前端函数或事件：检查更新按钮、checkUpdate、Check、OpenDownload
API/后端入口：GET /desktop-version.json
后端核心文件：/server/static/desktop-version.json
数据影响：无数据库写入，读取静态版本信息
状态变化：发现服务器版本高于当前版本时弹窗提示，用户确认后打开服务器下载地址；手动检查无新版时提示已是最新版
关联功能：DESKTOP_TRAY_001、AUTH_LOGIN_001
风险等级：中
验证方式：提高 desktop-version.json 中 version 后启动客户端或点击检查更新按钮，应弹出更新提示并能打开 downloadUrl；版本相同时手动检查提示已是最新版
