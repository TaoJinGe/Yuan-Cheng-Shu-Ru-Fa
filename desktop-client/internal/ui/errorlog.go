package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lxn/walk"
)

func WriteStartupError(err error) {
	path := errorLogPath()
	message := fmt.Sprintf("%s\n%v\n", time.Now().Format(time.RFC3339), err)
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	_ = os.WriteFile(path, []byte(message), 0600)
	_ = walk.MsgBox(nil, "远程语音输入启动失败", message+"\n日志："+path, walk.MsgBoxIconError)
}

func errorLogPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "voice-bridge-client-error.log"
	}
	return filepath.Join(dir, "voice-bridge-client", "error.log")
}
