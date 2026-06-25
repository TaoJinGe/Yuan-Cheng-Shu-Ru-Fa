package records

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode"
)

type Recorder struct {
	root string
	mu   sync.Mutex
}

func New(root string) *Recorder {
	return &Recorder{root: root}
}

func (r *Recorder) Save(username, text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	dir := filepath.Join(r.root, safeName(username))
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	now := time.Now()
	name := now.Format("2006-01-02") + ".md"
	path := filepath.Join(dir, name)
	entry := fmt.Sprintf("## %s\n\n%s\n\n---\n\n", now.Format("15:04:05"), text)

	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		if _, err := fmt.Fprintf(file, "# %s\n\n", now.Format("2006-01-02")); err != nil {
			return err
		}
	}
	_, err = file.WriteString(entry)
	return err
}

func safeName(value string) string {
	var b strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}
