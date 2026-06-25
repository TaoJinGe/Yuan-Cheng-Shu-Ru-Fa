//go:build windows

package config

import (
	"encoding/base64"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	cryptProtectData   = crypt32.NewProc("CryptProtectData")
	cryptUnprotectData = crypt32.NewProc("CryptUnprotectData")
	localFree          = kernel32.NewProc("LocalFree")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

func encryptLocal(text string) (string, error) {
	data := []byte(text)
	in := dataBlob{cbData: uint32(len(data)), pbData: &data[0]}
	var out dataBlob
	ok, _, err := cryptProtectData.Call(
		uintptr(unsafe.Pointer(&in)), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&out)),
	)
	if ok == 0 {
		return "", err
	}
	defer localFree.Call(uintptr(unsafe.Pointer(out.pbData)))
	bytes := unsafe.Slice(out.pbData, out.cbData)
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func decryptLocal(secret string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(secret)
	if err != nil || len(data) == 0 {
		return "", fmt.Errorf("bad password secret")
	}
	in := dataBlob{cbData: uint32(len(data)), pbData: &data[0]}
	var out dataBlob
	ok, _, callErr := cryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&in)), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&out)),
	)
	if ok == 0 {
		return "", callErr
	}
	defer localFree.Call(uintptr(unsafe.Pointer(out.pbData)))
	bytes := unsafe.Slice(out.pbData, out.cbData)
	return string(bytes), nil
}
