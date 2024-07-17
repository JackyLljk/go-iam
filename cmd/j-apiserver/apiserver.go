package main

import (
	"j-iam/internal/apiserver"
	"math/rand"
	"time"
)

func main() {
	// 初始化随机数生成器的种子，确保整个程序中的随机数生成器都使用这个种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 1. 调用 apiserver.NewApp 构建应用框架
	// 2. 执行 a.Run() 启动应用
	apiserver.NewApp("j-apiserver").Run()
}
