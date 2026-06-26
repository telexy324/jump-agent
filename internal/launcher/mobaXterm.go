package launcher

import (
	"fmt"
	"jump-agent/internal/model"
	"log"
	"os/exec"
	"sync/atomic"
	"time"
)

type MobaXterm struct{}

var crtStarted atomic.Bool

func ensureSecureMoba(path string) {
	if crtStarted.Load() {
		return
	}

	exec.Command(path).Start()
	time.Sleep(3000 * time.Millisecond) // 非常关键
	crtStarted.Store(true)
}

func (m *MobaXterm) Launch(c *model.SessionPayload) error {
	path, err := detectOrAsk("MobaXterm", findDefaultMobaXterm())
	if err != nil {
		return err
	}

	sshCmd := fmt.Sprintf(
		"ssh %s@%s -p %d",
		c.Secret,
		c.BastionHost,
		c.BastionPort,
	)

	ensureSecureMoba(path)

	args := []string{
		"-newtab",
		sshCmd,
	}

	log.Printf("Exec: %s %v", path, args)
	return exec.Command(path, args...).Start()
}

func findDefaultMobaXterm() []string {
	return []string{
		//`E:\SecureCRT\SecureCRT\SecureCRT.exe`,
	}
}
