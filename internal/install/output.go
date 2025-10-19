package install

import (
	"fmt"
	"path/filepath"
)

func OutputLog(homeDir, homeFile string, status *Status) {
	var resultStr string
	var colorCode string

	switch status.Result {
	case AlreadyLinked:
		resultStr = "already linked:"
		colorCode = "\033[90m" // gray
	case LinkUpdated:
		resultStr = "link updated:"
		colorCode = "\033[33m" // yellow
	case LinkCreated:
		resultStr = "link created:"
		colorCode = "\033[32m" // green
	default:
		resultStr = "result:"
		colorCode = ""
	}

	relPath, err := filepath.Rel(homeDir, homeFile)
	if err != nil {
		relPath = homeFile
	}

	msg := ""
	if colorCode != "" {
		msg += colorCode
	}
	msg += resultStr
	if colorCode != "" {
		msg += "\033[0m"
	}
	msg += " " + relPath
	if status.Backuped {
		msg += " (backup)"
	}
	fmt.Println(msg)
}
