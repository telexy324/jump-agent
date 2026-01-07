package model

type ConnInfo struct {
	JumpHost string `json:"jump_host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Protocol string `json:"protocol"` // ssh / sftp
	Client   string `json:"client"`   // securecrt / filezilla
	Password string `json:"password"` // 临时凭证（一次性）
}
