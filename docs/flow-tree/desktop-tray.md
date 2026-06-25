# desktop-tray FLOW_TREE

FLOW_ID：DESKTOP_TRAY_001
功能名称：Windows 客户端托盘常驻
入口页面：Windows 客户端登录窗口
触发点：点击窗口关闭按钮、托盘右键菜单
前端文件：/desktop-client/main.go、/desktop-client/internal/singleinstance/single_instance_windows.go、/desktop-client/internal/ui/run.go、/desktop-client/internal/ui/tray.go
前端函数或事件：main、Acquire、ShowAlreadyRunning、afterLogin、Closing、setupTray
API/后端入口：无
后端核心文件：无
数据影响：无
状态变化：双击启动先抢占单实例锁，已有客户端运行时只弹出重复启动提示并退出；首次启动显示窗口，登录成功后窗口继续显示，关闭窗口时隐藏到托盘，托盘退出时停止后台连接并退出进程
关联功能：AUTH_LOGIN_001、PRESENCE_SYNC_001
风险等级：中
验证方式：双击 exe 出现 GUI，再次双击只出现“软件已经在运行”提示且不会新增托盘图标；登录成功后窗口仍显示状态，点叉号后窗口隐藏且右下角有图标，托盘菜单可显示窗口和退出软件
