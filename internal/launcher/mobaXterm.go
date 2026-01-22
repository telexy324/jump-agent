package launcher

import (
	"fmt"
	"jump-agent/internal/model"
	"os/exec"
)

type MobaXterm struct{}

//func (s *SecureCRT) Launch(c *model.ConnInfo) error {
//	path, err := config.GetSecureCRTPath()
//	if err != nil {
//		return err
//	}
//
//	args := []string{
//		"/SSH2",
//		"/L", c.User,
//		"/P", strconv.Itoa(c.Port),
//		c.JumpHost,
//	}
//
//	return exec.Command(path, args...).Start()
//}

func (m *MobaXterm) Launch(c *model.SessionPayload) error {
	path, err := detectOrAsk("MobaXterm", findDefaultMobaXterm())
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf(
		"%s -ssh %s@%s -P %d -t target-host",
		path,
		c.Secret,
		c.BastionHost,
		c.BastionPort,
	)
	return exec.Command("cmd", "/C", cmd).Start()
}

func findDefaultMobaXterm() []string {
	return []string{
		`C:\Program Files\VanDyke Software\SecureCRT\SecureCRT.exe`,
		`C:\Program Files (x86)\VanDyke Software\SecureCRT\SecureCRT.exe`,
		//`E:\SecureCRT\SecureCRT\SecureCRT.exe`,
	}
}
