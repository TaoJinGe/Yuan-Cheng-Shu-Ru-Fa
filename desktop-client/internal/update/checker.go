package update

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const CurrentVersion = "1.1.7"

type Info struct {
	Version     string `json:"version"`
	DownloadURL string `json:"downloadUrl"`
	Notes       string `json:"notes"`
}

func Check(serverURL string) (Info, bool, error) {
	checkURL := versionURL(serverURL)
	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(checkURL)
	if err != nil {
		return Info{}, false, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Info{}, false, nil
	}
	var info Info
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return Info{}, false, err
	}
	return info, newer(info.Version, CurrentVersion), nil
}

func OpenDownload(rawURL string) error {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL).Start()
}

func versionURL(serverURL string) string {
	u, err := url.Parse(serverURL)
	if err != nil {
		return serverURL
	}
	if u.Scheme == "wss" {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Path = "/desktop-version.json"
	u.RawQuery = ""
	return u.String()
}

func newer(remote, current string) bool {
	left := splitVersion(remote)
	right := splitVersion(current)
	for i := 0; i < len(left) || i < len(right); i++ {
		a, b := part(left, i), part(right, i)
		if a != b {
			return a > b
		}
	}
	return false
}

func splitVersion(value string) []string {
	return strings.Split(strings.TrimPrefix(value, "v"), ".")
}

func part(items []string, index int) int {
	if index >= len(items) {
		return 0
	}
	value, _ := strconv.Atoi(items[index])
	return value
}
