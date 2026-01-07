package main

import (
	"log"
	"os"

	"jump-agent/internal/protocol"
	"jump-agent/internal/registry"
)

func main() {
	// 1️⃣ 确保自定义协议已注册
	if err := registry.EnsureMyJumpProtocol(); err != nil {
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
