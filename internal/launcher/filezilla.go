package launcher

import (
	"fmt"
	"jump-agent/internal/model"
	"os/exec"
)

type FileZilla struct{}

//func (f *FileZilla) Launch(c *model.ConnInfo) error {
//	path, err := config.GetFileZillaPath()
//	if err != nil {
//		return err
//	}
//
//	// FileZilla 支持 sftp:// URL
//	// sftp://user:pass@host:port/
//	url := fmt.Sprintf(
//		"sftp://%s:%s@%s:%d/",
//		c.User,
//		c.Password,
//		c.JumpHost,
//		c.Port,
//	)
//
//	return exec.Command(path, url).Start()
//}

func (f *FileZilla) Launch(c *model.SessionPayload) error {
	path, err := detectOrAsk("FileZilla", findDefaultFileZilla())
	if err != nil {
		return err
	}
	url := fmt.Sprintf(
		"sftp://%s@%s:%d/",
		c.Secret,
		c.BastionHost,
		c.BastionPort,
	)
	cmd := exec.Command(path, url)
	return cmd.Start()
}

func findDefaultFileZilla() []string {
	return []string{
		`C:\Program Files\FileZilla FTP Client\filezilla.exe`,
		`C:\Program Files (x86)\FileZilla FTP Client\filezilla.exe`,
	}
}
