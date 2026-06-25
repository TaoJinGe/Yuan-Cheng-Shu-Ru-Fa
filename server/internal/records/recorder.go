package records

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type Recorder struct {
	root string
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
	name := now.Format("2006-01-02_15-04-05.000") + ".md"
	content := fmt.Sprintf("# %s\n\n%s\n", now.Format(time.RFC3339), text)
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0600)
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
