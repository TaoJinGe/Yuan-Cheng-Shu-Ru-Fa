# auth FLOW_TREE

FLOW_ID：AUTH_REGISTER_001
功能名称：用户注册账号并获取 token
入口页面：手机网页登录页 /，Windows 客户端登录窗口
触发点：注册按钮
前端文件：/server/static/app.js、/desktop-client/internal/ui/run.go、/desktop-client/internal/client/api.go
前端函数或事件：手机端 registerAccount，桌面端 register
API/后端入口：POST /api/register
后端核心文件：/server/internal/web/handlers.go、/server/internal/auth/store.go、/server/internal/auth/users.go
数据影响：SQLite app.db users 表、服务端内存 sessions
状态变化：保存 bcrypt 密码哈希，生成 24 小时或 30 天 token
关联功能：AUTH_LOGIN_001、PRESENCE_SYNC_001
风险等级：高
验证方式：新账号注册成功后手机端和桌面端可用同账号登录

FLOW_ID：AUTH_LOGIN_001
功能名称：用户登录并获取 token
入口页面：手机网页登录页 /，Windows 客户端登录窗口
触发点：登录按钮或客户端 token 自动登录
前端文件：/server/static/app.js、/desktop-client/internal/ui/run.go、/desktop-client/internal/client/api.go
前端函数或事件：手机端 login，桌面端 login
API/后端入口：POST /api/login
后端核心文件：/server/internal/web/handlers.go、/server/internal/auth/store.go、/server/internal/auth/users.go
数据影响：SQLite app.db users 表只读、服务端内存 sessions
状态变化：生成 24 小时或 30 天 token，客户端保存 token
关联功能：AUTH_REGISTER_001、PRESENCE_SYNC_001、TEXT_SEND_001
风险等级：中
验证方式：勾选和不勾选 30 天免登录后检查 token 过期时间和本地保存行为

FLOW_ID：AUTH_LOGOUT_001
功能名称：用户退出登录并清理 token
入口页面：手机网页 /，Windows 客户端登录窗口
触发点：退出登录按钮
前端文件：/server/static/app.js、/desktop-client/internal/ui/run.go、/desktop-client/internal/client/api.go
前端函数或事件：手机端 logout click，桌面端 logout
API/后端入口：POST /api/logout
后端核心文件：/server/internal/web/handlers.go、/server/internal/auth/store.go
数据影响：服务端内存 sessions
状态变化：吊销 token，本地删除 token
关联功能：PRESENCE_SYNC_001
风险等级：中
验证方式：退出后刷新网页或重启客户端应回到未登录状态
