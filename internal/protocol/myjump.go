package protocol

import (
	"jump-agent/internal/launcher"
	"jump-agent/internal/token"
	"strings"
)

//func Handle(raw string) error {
//	u, err := url.Parse(raw)
//	if err != nil {
//		return err
//	}
//
//	if u.Scheme != "myjump" {
//		return fmt.Errorf("invalid scheme")
//	}
//
//	tokenStr := u.Query().Get("token")
//	if tokenStr == "" {
//		return fmt.Errorf("token missing")
//	}
//
//	// 1️⃣ 向堡垒机校验并消费 token
//	conn, err := token.Consume(tokenStr)
//	if err != nil {
//		return err
//	}
//
//	// 2️⃣ 启动 SecureCRT
//	return launcher.Default().Launch(conn)
//}

func Handle(raw string) error {
	//u, err := url.Parse(raw)
	//if err != nil {
	//	return err
	//}
	//
	//tokenStr := u.Query().Get("token")
	//if tokenStr == "" {
	//	return fmt.Errorf("token missing")
	//}
	tokenStr := strings.TrimPrefix(raw, "myjump://")
	tokenStr = strings.TrimSuffix(tokenStr, "/")
	conn, err := token.Consume(tokenStr)
	if err != nil {
		return err
	}

	return launcher.Get(conn.Client).Launch(conn)
}
