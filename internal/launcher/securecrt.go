package launcher

import (
	"jump-agent/internal/model"
	"log"
	"os/exec"
	"strconv"
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

//var crtStarted atomic.Bool
//
//func ensureSecureCRT(path string) {
//	if crtStarted.Load() {
//		return
//	}
//
//	exec.Command(path).Start()
//	time.Sleep(1500 * time.Millisecond) // 非常关键
//	crtStarted.Store(true)
//}

func (s *SecureCRT) Launch(c *model.SessionPayload) error {
	path, err := detectOrAsk("SecureCRT", findDefaultSecureCRT())
	if err != nil {
		return err
	}

	//ensureSecureCRT(path)

	args := []string{
		"/T",
		"/SSH2",
		"/L", c.Secret,
		"/P", strconv.Itoa(c.BastionPort),
		c.BastionHost,
	}
	cmd := exec.Command(path, args...)
	log.Printf("Exec: %s %v\n", path, args)
	return cmd.Start()
}

func findDefaultSecureCRT() []string {
	return []string{
		//`C:\Program Files\VanDyke Software\SecureCRT\SecureCRT.exe`,
		//`C:\Program Files (x86)\VanDyke Software\SecureCRT\SecureCRT.exe`,
		//`E:\SecureCRT\SecureCRT\SecureCRT.exe`,
	}
}
