//go:build windows

package paste

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	cfUnicodeText = 13
	gmemMoveable  = 0x0002
	vkControl     = 0x11
	vkV           = 0x56
	keyEventUp    = 0x0002
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	openClipboard            = user32.NewProc("OpenClipboard")
	closeClipboard           = user32.NewProc("CloseClipboard")
	emptyClipboard           = user32.NewProc("EmptyClipboard")
	getClipboardData         = user32.NewProc("GetClipboardData")
	setClipboardData         = user32.NewProc("SetClipboardData")
	keybdEvent               = user32.NewProc("keybd_event")
	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessID = user32.NewProc("GetWindowThreadProcessId")
	getGUIThreadInfo         = user32.NewProc("GetGUIThreadInfo")
	globalAlloc              = kernel32.NewProc("GlobalAlloc")
	globalLock               = kernel32.NewProc("GlobalLock")
	globalUnlock             = kernel32.NewProc("GlobalUnlock")
	globalSize               = kernel32.NewProc("GlobalSize")
)

type rect struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type guiThreadInfo struct {
	CbSize        uint32
	Flags         uint32
	HwndActive    uintptr
	HwndFocus     uintptr
	HwndCapture   uintptr
	HwndMenuOwner uintptr
	HwndMoveSize  uintptr
	HwndCaret     uintptr
	RcCaret       rect
}

func PasteText(text string, restoreDelay time.Duration) error {
	if !foregroundWindowExists() {
		return fmt.Errorf("当前电脑没有活动窗口，已忽略")
	}
	previous, _ := readClipboardText()
	if err := writeClipboardText(text); err != nil {
		return err
	}
	time.Sleep(120 * time.Millisecond)
	sendCtrlV()
	if previous != "" && restoreDelay > 0 {
		time.Sleep(restoreDelay)
		_ = writeClipboardText(previous)
	}
	return nil
}

func readClipboardText() (string, error) {
	if !open() {
		return "", fmt.Errorf("无法打开剪贴板")
	}
	defer closeClipboard.Call()

	handle, _, _ := getClipboardData.Call(cfUnicodeText)
	if handle == 0 {
		return "", nil
	}
	ptr, _, _ := globalLock.Call(handle)
	if ptr == 0 {
		return "", fmt.Errorf("无法读取剪贴板")
	}
	defer globalUnlock.Call(handle)

	size, _, _ := globalSize.Call(handle)
	if size == 0 {
		return "", nil
	}
	words := unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), int(size)/2)
	n := 0
	for n < len(words) && words[n] != 0 {
		n++
	}
	return syscall.UTF16ToString(words[:n]), nil
}

func writeClipboardText(text string) error {
	data := syscall.StringToUTF16(text)
	bytes := uintptr(len(data) * 2)
	handle, _, _ := globalAlloc.Call(gmemMoveable, bytes)
	if handle == 0 {
		return fmt.Errorf("无法分配剪贴板内存")
	}
	ptr, _, _ := globalLock.Call(handle)
	if ptr == 0 {
		return fmt.Errorf("无法锁定剪贴板内存")
	}
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), len(data)), data)
	globalUnlock.Call(handle)

	if !open() {
		return fmt.Errorf("无法打开剪贴板")
	}
	defer closeClipboard.Call()
	emptyClipboard.Call()
	result, _, _ := setClipboardData.Call(cfUnicodeText, handle)
	if result == 0 {
		return fmt.Errorf("无法写入剪贴板")
	}
	return nil
}

func open() bool {
	for i := 0; i < 8; i++ {
		ok, _, _ := openClipboard.Call(0)
		if ok != 0 {
			return true
		}
		time.Sleep(25 * time.Millisecond)
	}
	return false
}

func sendCtrlV() {
	keybdEvent.Call(vkControl, 0, 0, 0)
	keybdEvent.Call(vkV, 0, 0, 0)
	keybdEvent.Call(vkV, 0, keyEventUp, 0)
	keybdEvent.Call(vkControl, 0, keyEventUp, 0)
}

func foregroundWindowExists() bool {
	fg, _, _ := getForegroundWindow.Call()
	return fg != 0
}

func inputCaretActive() bool {
	fg, _, _ := getForegroundWindow.Call()
	if fg == 0 {
		return false
	}
	threadID, _, _ := getWindowThreadProcessID.Call(fg, 0)
	if threadID == 0 {
		return false
	}
	info := guiThreadInfo{CbSize: uint32(unsafe.Sizeof(guiThreadInfo{}))}
	ok, _, _ := getGUIThreadInfo.Call(threadID, uintptr(unsafe.Pointer(&info)))
	return ok != 0 && info.HwndCaret != 0
}
