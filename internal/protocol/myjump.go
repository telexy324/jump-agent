package protocol

import (
	"errors"
	"jump-agent/internal/launcher"
	"jump-agent/internal/model"
	"strings"
	"time"
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
	//conn, err := token.Consume(tokenStr)
	//if err != nil {
	//	return err
	//}
	//conn := model.ConnInfo{
	//	JumpHost: "",
	//	Port:     0,
	//	User:     "",
	//	Protocol: "",
	//	Client:   "",
	//	Password: "",
	//}
	conns, err := model.ParseSession(tokenStr)
	if err != nil {
		return err
	}
	if len(conns) == 0 {
		return errors.New("invalid token, no connections")
	}

	if len(conns) == 1 {
		return launcher.Get(conns[0].Client).Launch(conns[0])
	}
	if err := launcher.Get(conns[0].Client).Launch(conns[0]); err != nil {
		return err
	}
	time.Sleep(2000 * time.Millisecond)
	for _, conn := range conns[1:] {
		if err := launcher.Get(conn.Client).Launch(conn); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}
