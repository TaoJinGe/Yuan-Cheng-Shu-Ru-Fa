# server TREE_MAP

main.go —— 服务端启动入口，只负责依赖装配和监听启动（修改影响服务进程启动）
go.mod —— 服务端模块依赖文件，只负责声明构建依赖（修改影响依赖解析）
go.sum —— 服务端依赖校验文件，只负责锁定依赖校验值（修改影响依赖一致性）
internal/config/ —— 配置模块目录，只负责服务端端口、SQLite、记录目录和 token 期限配置（修改影响启动参数）
internal/config/config.go —— 配置读取模块，只负责解析 flag 与环境变量（修改影响配置优先级）
internal/auth/ —— 登录会话模块目录，只负责账号认证和 token 生命周期管理（修改影响认证状态，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、AUTH_LOGOUT_001）
internal/auth/store.go —— 认证编排模块，只负责注册、登录、token 签发、校验和吊销（修改影响登录有效期）
internal/auth/users.go —— SQLite 用户存储模块，只负责 app.db 用户表初始化、写入和账号查询（修改影响账号持久化，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001）
internal/records/ —— 文本记录目录，只负责按账号和时间保存手机发送内容（修改影响发送记录）
internal/records/recorder.go —— 文本记录模块，只负责 records/账号/时间.md 文件写入（修改影响发送记录，关联 FLOW_ID：TEXT_SEND_001）
internal/hub/ —— WebSocket Hub 模块目录，只负责在线连接集合、presence 广播和文本转发（修改影响状态灯和发送链路，关联 FLOW_ID：PRESENCE_SYNC_001、TEXT_SEND_001）
internal/hub/message.go —— WebSocket 消息结构定义，只负责端到端消息协议（修改影响前后端消息兼容）
internal/hub/hub.go —— 在线连接 Hub，只负责连接注册、断开、状态计算和房间转发（修改影响实时状态）
internal/web/ —— HTTP 处理模块目录，只负责 API、静态文件和 WebSocket 升级入口（修改影响服务端路由）
internal/web/handlers.go —— HTTP 路由处理模块，只负责注册、登录、退出、健康检查、静态页和 WebSocket 入口（修改影响 API 行为，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、AUTH_LOGOUT_001、PRESENCE_SYNC_001）
static/ —— 手机网页静态资源目录，只负责移动端页面资源（修改影响手机端体验）
static/index.html —— 手机端 HTML 结构，只负责登录页和无配对码输入页两态容器（修改影响移动端布局，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、TEXT_SEND_001）
static/favicon.ico —— 手机网页图标文件，只负责浏览器标签页和避免 favicon 404（修改影响静态资源请求）
static/app.js —— 手机端交互脚本，只负责注册登录、页面切换、WebSocket、状态灯和同账号发送按钮（修改影响手机端交互，关联 FLOW_ID：AUTH_REGISTER_001、AUTH_LOGIN_001、PRESENCE_SYNC_001、TEXT_SEND_001）
static/style.css —— 手机端样式文件，只负责登录页、大输入框、发送按钮和状态灯样式（修改影响移动端视觉）
