package launcher

import "jump-agent/internal/model"

type Launcher interface {
	Launch(*model.ConnInfo) error
}

func Get(name string) Launcher {
	switch name {
	case "filezilla":
		return &FileZilla{}
	case "securecrt":
		fallthrough
	default:
		return &SecureCRT{}
	}
}
