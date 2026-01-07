package launcher

import (
	"os/exec"
	"strconv"

	"jump-agent/internal/config"
	"jump-agent/internal/model"
)

type SecureCRT struct{}

func (s *SecureCRT) Launch(c *model.ConnInfo) error {
	path, err := config.GetSecureCRTPath()
	if err != nil {
		return err
	}

	args := []string{
		"/SSH2",
		"/L", c.User,
		"/P", strconv.Itoa(c.Port),
		c.JumpHost,
	}

	return exec.Command(path, args...).Start()
}
