//go:build windows

package singleinstance

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	user32            = syscall.NewLazyDLL("user32.dll")
	createMutexW      = kernel32.NewProc("CreateMutexW")
	messageBoxW       = user32.NewProc("MessageBoxW")
	ErrAlreadyRunning = errors.New("voice bridge client is already running")
)

const (
	errorAlreadyExists = syscall.Errno(183)
	messageBoxIconInfo = 0x40
)

type Lock struct {
	handle syscall.Handle
}

func Acquire(name string) (*Lock, error) {
	mutexName, err := syscall.UTF16PtrFromString("Global\\" + name)
	if err != nil {
		return nil, err
	}
	handle, _, callErr := createMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	if handle == 0 {
		return nil, callErr
	}
	if callErr == errorAlreadyExists {
		_ = syscall.CloseHandle(syscall.Handle(handle))
		return nil, ErrAlreadyRunning
	}
	return &Lock{handle: syscall.Handle(handle)}, nil
}

func (l *Lock) Release() {
	if l == nil || l.handle == 0 {
		return
	}
	_ = syscall.CloseHandle(l.handle)
	l.handle = 0
}

func ShowAlreadyRunning() {
	title, _ := syscall.UTF16PtrFromString("远程语音输入")
	text, _ := syscall.UTF16PtrFromString("软件已经在运行了。\n\n请在右下角托盘图标中显示窗口，或先退出已有软件后再重新打开。")
	messageBoxW.Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), messageBoxIconInfo)
}
