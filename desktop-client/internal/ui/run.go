package ui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"voice-bridge-client/internal/app"
	"voice-bridge-client/internal/client"
)

type windowState struct {
	app        *app.State
	mw         *walk.MainWindow
	notifyIcon *walk.NotifyIcon
	serverEdit *walk.LineEdit
	userEdit   *walk.LineEdit
	passEdit   *walk.LineEdit
	remember   *walk.CheckBox
	status     *walk.Label
	loginBtn   *walk.PushButton
}

func Run(defaultServer string) error {
	state, err := app.New(defaultServer)
	if err != nil {
		return err
	}
	ws := &windowState{app: state}
	if err := ws.build(); err != nil {
		return err
	}
	if err := ws.setupTray(); err != nil {
		return err
	}
	ws.loadConfig()
	ws.mw.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		ws.mw.Hide()
	})
	if state.Config().Valid() {
		ws.afterLogin()
	}
	ws.mw.Run()
	ws.cleanup()
	return nil
}

func (w *windowState) build() error {
	if err := (MainWindow{
		AssignTo: &w.mw,
		Title:    "远程语音输入",
		Size:     Size{Width: 420, Height: 360},
		MinSize:  Size{Width: 390, Height: 340},
		Font:     Font{Family: "Microsoft YaHei UI", PointSize: 10},
		Layout:   VBox{Margins: Margins{Left: 18, Top: 16, Right: 18, Bottom: 16}, Spacing: 9},
		Children: []Widget{
			Label{
				Text:      "远程语音输入",
				Font:      Font{Family: "Microsoft YaHei UI", PointSize: 13, Bold: true},
				MinSize:   Size{Height: 24},
				TextColor: walk.RGB(23, 32, 38),
			},
			Label{Text: "服务器 WebSocket 地址", TextColor: walk.RGB(80, 96, 106)},
			LineEdit{
				AssignTo:  &w.serverEdit,
				CueBanner: "wss://srf.cccz.cc/ws",
				MinSize:   Size{Height: 34},
			},
			Label{Text: "账号", TextColor: walk.RGB(80, 96, 106)},
			LineEdit{
				AssignTo:  &w.userEdit,
				CueBanner: "请输入账号",
				MinSize:   Size{Height: 34},
			},
			Label{Text: "密码", TextColor: walk.RGB(80, 96, 106)},
			LineEdit{
				AssignTo:     &w.passEdit,
				PasswordMode: true,
				CueBanner:    "请输入密码",
				MinSize:      Size{Height: 34},
			},
			CheckBox{AssignTo: &w.remember, Text: "30 天免登录", Checked: true, MinSize: Size{Height: 28}},
			Composite{
				Layout: HBox{MarginsZero: true, Spacing: 8},
				Children: []Widget{
					PushButton{
						AssignTo:  &w.loginBtn,
						Text:      "登录",
						MinSize:   Size{Height: 38},
						OnClicked: w.login,
					},
					PushButton{Text: "注册", MinSize: Size{Height: 38}, OnClicked: w.register},
					PushButton{Text: "退出登录", MinSize: Size{Height: 38}, OnClicked: w.logout},
				},
			},
			Label{
				AssignTo:  &w.status,
				Text:      "未登录",
				MinSize:   Size{Height: 30},
				TextColor: walk.RGB(20, 108, 92),
			},
		},
	}).Create(); err != nil {
		return err
	}
	if icon, err := loadAppIcon(); err == nil {
		_ = w.mw.SetIcon(icon)
	}
	return nil
}

func (w *windowState) loadConfig() {
	cfg := w.app.Config()
	w.serverEdit.SetText(cfg.ServerURL)
	w.userEdit.SetText(cfg.UserID)
}

func (w *windowState) login() {
	w.doAuth(false)
}

func (w *windowState) register() {
	w.doAuth(true)
}

func (w *windowState) doAuth(register bool) {
	if register {
		w.setStatus("正在注册...")
	} else {
		w.setStatus("正在登录...")
	}
	input := app.LoginInput{
		ServerURL: w.serverEdit.Text(),
		Username:  w.userEdit.Text(),
		Password:  w.passEdit.Text(),
		Remember:  w.remember.Checked(),
	}
	var err error
	if register {
		err = w.app.Register(input)
	} else {
		err = w.app.Login(input)
	}
	if err != nil {
		w.setStatus("登录失败：" + err.Error())
		return
	}
	w.passEdit.SetText("")
	w.afterLogin()
}

func (w *windowState) afterLogin() {
	w.setStatus("服务器连接中")
	w.app.StartSocket(client.Events{
		Status:     w.setStatus,
		Connected:  w.setConnected,
		AuthFailed: w.authFailed,
	})
}

func (w *windowState) logout() {
	w.app.Logout()
	w.setStatus("未登录")
	w.mw.Show()
}

func (w *windowState) setStatus(text string) {
	w.mw.Synchronize(func() {
		w.status.SetText(text)
		if w.notifyIcon != nil {
			w.notifyIcon.SetToolTip("远程语音输入 - " + text)
		}
	})
}

func (w *windowState) setConnected(ok bool) {
	if ok {
		w.notify("手机端已连接")
	}
}

func (w *windowState) authFailed() {
	w.mw.Synchronize(func() {
		w.app.Logout()
		w.status.SetText("登录已过期，请重新登录")
		w.mw.Show()
	})
}

func (w *windowState) notify(text string) {
	if w.notifyIcon != nil {
		_ = w.notifyIcon.ShowInfo("远程语音输入", text)
	}
}

func (w *windowState) cleanup() {
	w.app.StopSocket()
	if w.notifyIcon != nil {
		w.notifyIcon.Dispose()
	}
}
