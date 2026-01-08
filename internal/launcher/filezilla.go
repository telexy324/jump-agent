package launcher

import (
	"fmt"
	"jump-agent/internal/config"
	"jump-agent/internal/model"
	"os/exec"
)

type FileZilla struct{}

func (f *FileZilla) Launch(c *model.ConnInfo) error {
	path, err := config.GetFileZillaPath()
	if err != nil {
		return err
	}

	// FileZilla 支持 sftp:// URL
	// sftp://user:pass@host:port/
	url := fmt.Sprintf(
		"sftp://%s:%s@%s:%d/",
		c.User,
		c.Password,
		c.JumpHost,
		c.Port,
	)

	return exec.Command(path, url).Start()
}
