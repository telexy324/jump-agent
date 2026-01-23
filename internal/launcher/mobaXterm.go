package launcher

import (
	"fmt"
	"jump-agent/internal/model"
	"log"
	"os/exec"
	"strconv"
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

	args := []string{
		"-ssh",
		fmt.Sprintf("%s@%s", c.Secret, c.BastionHost),
		"-P",
		strconv.Itoa(c.BastionPort),
	}

	log.Printf("Exec: %s %v", path, args)
	return exec.Command(path, args...).Start()
}

func findDefaultMobaXterm() []string {
	return []string{
		`C:\Program Files\VanDyke Software\SecureCRT\SecureCRT.exe`,
		`C:\Program Files (x86)\VanDyke Software\SecureCRT\SecureCRT.exe`,
		//`E:\SecureCRT\SecureCRT\SecureCRT.exe`,
	}
}
