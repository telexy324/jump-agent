package launcher

//func detectOrAsk(app string, candidates []string) (string, error) {
//	// 1. 读取注册表缓存
//	if p, err := agent.LoadPath(app); err == nil {
//		if _, err := os.Stat(p); err == nil {
//			return p, nil
//		}
//	}
//
//	// 2. 尝试默认路径
//	for _, c := range candidates {
//		if _, err := os.Stat(c); err == nil {
//			agent.SavePath(app, c)
//			return c, nil
//		}
//	}
//
//	// 3. 用户选择 EXE
//	exe, err := agent.SelectExecutable()
//	if err != nil {
//		return "", errors.New("用户未选择可执行文件")
//	}
//
//	agent.SavePath(app, exe)
//	return exe, nil
//}

//func LaunchSecureCRT(args []string) error {
//	path, err := detectOrAsk("SecureCRT", findDefaultSecureCRT())
//	if err != nil {
//		return err
//	}
//	cmd := exec.Command(path, args...)
//	return cmd.Start()
//}
//
//func LaunchFileZilla(args []string) error {
//	path, err := detectOrAsk("FileZilla", findDefaultFileZilla())
//	if err != nil {
//		return err
//	}
//	cmd := exec.Command(path, args...)
//	return cmd.Start()
//}
