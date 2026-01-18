package registry

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func EnsureMyJumpProtocol() error {
	exe, _ := os.Executable()
	exe, _ = filepath.Abs(exe)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		"myjump",
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	key.SetStringValue("", "URL:MyJump Protocol")
	key.SetStringValue("URL Protocol", "")

	cmdKey, _, err := registry.CreateKey(
		key,
		"shell\\open\\command",
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer cmdKey.Close()

	cmd := fmt.Sprintf("\"%s\" \"%%1\"", exe)
	return cmdKey.SetStringValue("", cmd)
}
