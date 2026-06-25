//go:build !windows

package paste

import (
	"fmt"
	"time"
)

func PasteText(_ string, _ time.Duration) error {
	return fmt.Errorf("第一版只支持 Windows 客户端")
}
