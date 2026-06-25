# ROOT_MAP_1

/README.md —— 项目使用说明入口，说明部署、客户端启动、手机网页使用和常见问题（修改影响文档入口）
/.gitignore —— Git 忽略规则文件，屏蔽运行数据、构建产物和临时文件（修改影响仓库上传范围）
/启动本地服务.cmd —— Windows 本地启动脚本，负责同步手机端静态资源、在 dist/server 工作目录启动服务端并打开默认浏览器（修改影响本地调试入口）
/关闭本地服务.cmd —— Windows 本地关闭脚本，负责关闭占用 8080 端口的本地服务（修改影响本地服务清理）
/server/ —— Go 服务端目录，负责 HTTP、WebSocket、会话、在线状态和手机静态页（修改影响服务端 API 与中转链路）
/server/main.go —— 服务端启动入口，只负责装配配置、会话、Hub 和 HTTP 路由（修改影响服务进程启动）
/server/go.mod —— 服务端 Go 模块定义，记录 WebSocket 依赖（修改影响服务端构建）
/server/go.sum —— 服务端依赖校验文件，记录模块版本校验值（修改影响依赖一致性）
/server/internal/config/ —— 服务端配置目录，负责命令行参数与环境变量读取（修改影响启动配置）
/server/internal/auth/ —— 服务端认证目录，负责 SQLite 账号注册、密码哈希、token 签发、校验和失效（修改影响登录态）
/server/internal/hub/ —— 服务端 WebSocket Hub 目录，负责同账号转发和设备在线状态（修改影响消息中转与状态灯）
/server/internal/records/ —— 服务端文本记录目录，负责按账号和时间保存手机发送内容（修改影响发送记录）
/server/internal/web/ —— 服务端 HTTP 目录，负责登录 API、健康检查、静态页和 WebSocket 入口（修改影响 HTTP/API 入口）
/server/static/ —— 手机网页静态资源目录，负责移动端登录、状态灯和发送交互（修改影响手机端 UI 与交互）
/desktop-client/ —— Windows 电脑客户端目录，负责登录、常驻 WebSocket、剪贴板写入和 Ctrl+V（修改影响桌面端体验）
/desktop-client/main.go —— 桌面客户端启动入口，只负责启动 Windows GUI 应用（修改影响客户端启动）
/desktop-client/go.mod —— 桌面客户端 Go 模块定义，记录 WebSocket 依赖（修改影响客户端构建）
/desktop-client/go.sum —— 桌面客户端依赖校验文件，记录模块版本校验值（修改影响依赖一致性）
/desktop-client/internal/config/ —— 桌面客户端配置目录，负责本地 token 配置读写（修改影响 30 天免登录）
/desktop-client/internal/client/ —— 桌面客户端连接目录，负责服务端登录与 WebSocket 消息处理（修改影响在线状态与收文）
/desktop-client/internal/app/ —— 桌面客户端应用状态目录，负责配置、登录状态和后台 WebSocket 生命周期（修改影响 GUI 行为）
/desktop-client/internal/ui/ —— 桌面客户端 Windows UI 目录，负责登录注册窗口和系统托盘（修改影响桌面端交互）
/desktop-client/internal/paste/ —— 桌面客户端粘贴目录，负责 Windows 剪贴板和 SendInput（修改影响自动粘贴）
/docs/AI_SYSTEM_RULES.md —— 项目唯一 AI 规则文件，定义索引、拆分、300 行限制和功能链路规则（修改影响 AI 执行约束）
/docs/root/ —— ROOT_MAP 目录，维护一级文件系统索引（修改影响文件定位）
/docs/tree/ —— TREE_MAP 目录，维护二级目录索引（修改影响目录内职责定位）
/docs/flow-root/ —— 功能链路总索引目录，维护功能域到 FLOW_TREE 的定位（修改影响功能链路定位）
/docs/flow-tree/ —— 功能链路详情目录，维护登录、状态、发送等端到端链路（修改影响功能追踪）
/dist/ —— 构建产物目录，保存 Linux 服务端、Windows 客户端和手机网页静态文件（生成产物，AI 不修改）
/data/ —— 本地运行数据目录，由启动脚本运行时生成 SQLite app.db 和 records 记录目录（运行时数据，AI 不预设内容）
