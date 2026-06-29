package protocol

import (
	"errors"
	"jump-agent/internal/launcher"
	"jump-agent/internal/model"
	"jump-agent/internal/tunnel"
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

	tunnels := make([]*tunnel.Session, 0, len(conns))
	localConns := make([]*model.SessionPayload, 0, len(conns))
	for _, conn := range conns {
		t, localConn, err := tunnel.Start(conn)
		if err != nil {
			for _, started := range tunnels {
				started.Close()
			}
			return err
		}
		tunnels = append(tunnels, t)
		localConns = append(localConns, localConn)
	}
	defer func() {
		for _, t := range tunnels {
			t.Close()
		}
	}()

	if len(conns) == 1 {
		if err := launcher.Get(localConns[0].Client).Launch(localConns[0]); err != nil {
			return err
		}
		tunnels[0].Wait()
		return nil
	}
	if err := launcher.Get(localConns[0].Client).Launch(localConns[0]); err != nil {
		return err
	}
	//time.Sleep(2000 * time.Millisecond)
	for _, conn := range localConns[1:] {
		if err := launcher.Get(conn.Client).Launch(conn); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	for _, t := range tunnels {
		t.Wait()
	}
	return nil
}
