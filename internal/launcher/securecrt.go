package launcher

import (
	"log"
	"os/exec"
	"strconv"

	"jump-agent/internal/model"
)

type SecureCRT struct{}

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

func (s *SecureCRT) Launch(c *model.ConnInfo) error {
	path, err := detectOrAsk("SecureCRT", findDefaultSecureCRT())
	if err != nil {
		return err
	}
	args := []string{
		"/SSH2",
		"/L", c.JumpHost,
		"/P", strconv.Itoa(c.Port),
		"/PASSWORD", c.Password,
		"127.0.0.1",
	}
	cmd := exec.Command(path, args...)
	log.Printf("Exec: %s %v\n", path, args)
	return cmd.Start()
}

func findDefaultSecureCRT() []string {
	return []string{
		`C:\Program Files\VanDyke Software\SecureCRT\SecureCRT.exe`,
		`C:\Program Files (x86)\VanDyke Software\SecureCRT\SecureCRT.exe`,
		`E:\SecureCRT\SecureCRT\SecureCRT.exe`,
	}
}
