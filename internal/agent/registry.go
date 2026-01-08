//go:build windows
// +build windows

package agent

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

const regKey = `Software\BastionAgent`

func SavePath(app string, exePath string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, regKey, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	return k.SetStringValue(app, exePath)
}

func LoadPath(app string) (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, regKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	r, _, e := k.GetStringValue(app)
	return r, e
}

func RegisterProtocol() error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Classes\bastion-ssh`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	k.SetStringValue("", "URL:Bastion SSH Protocol")
	k.SetStringValue("URL Protocol", "")

	// command
	cmdKey, _, _ := registry.CreateKey(k, `shell\open\command`, registry.ALL_ACCESS)
	defer cmdKey.Close()

	agentPath, _ := os.Executable()
	cmdKey.SetStringValue("", agentPath+" \"%1\"")
	return nil
}
