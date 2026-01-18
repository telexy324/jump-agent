package main

import (
	"jump-agent/internal/agent"
	"log"
	"os"
	"path/filepath"

	"jump-agent/internal/protocol"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(exePath)
	logPath := filepath.Join(dir, "myapp.log")

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	//mw := io.MultiWriter(os.Stdout, f)
	//log.SetOutput(mw)
	// 1️⃣ 确保自定义协议已注册
	if err := agent.RegisterProtocol(); err != nil {
		log.Fatal("register protocol failed:", err)
	}

	// 2️⃣ 没参数 = 常驻 or 配置模式
	if len(os.Args) < 2 {
		log.Println("agent running (no protocol invoke)")
		select {} // 后续可做托盘
	}

	// 3️⃣ 处理 myjump://
	url := os.Args[1]
	if err := protocol.Handle(url); err != nil {
		log.Println("handle protocol error:", err)
	}
}
