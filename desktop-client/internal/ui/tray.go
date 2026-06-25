package ui

import "github.com/lxn/walk"

func (w *windowState) setupTray() error {
	ni, err := walk.NewNotifyIcon(w.mw)
	if err != nil {
		return err
	}
	w.notifyIcon = ni
	icon, err := loadAppIcon()
	if err != nil {
		return err
	}
	if err := ni.SetIcon(icon); err != nil {
		return err
	}
	if err := ni.SetToolTip("远程语音输入"); err != nil {
		return err
	}
	showAction := walk.NewAction()
	showAction.SetText("显示窗口")
	showAction.Triggered().Attach(func() {
		w.mw.Show()
		w.mw.Activate()
	})
	exitAction := walk.NewAction()
	exitAction.SetText("退出软件")
	exitAction.Triggered().Attach(func() {
		w.app.StopSocket()
		w.notifyIcon.Dispose()
		w.mw.Dispose()
	})
	ni.ContextMenu().Actions().Add(showAction)
	ni.ContextMenu().Actions().Add(exitAction)
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			w.mw.Show()
			w.mw.Activate()
		}
	})
	return ni.SetVisible(true)
}
