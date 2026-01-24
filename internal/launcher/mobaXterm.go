package launcher

import (
	"fmt"
	"jump-agent/internal/model"
	"log"
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

	sshCmd := fmt.Sprintf(
		"ssh %s@%s -p %d",
		c.Secret,
		c.BastionHost,
		c.BastionPort,
	)

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
