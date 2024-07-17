package main

import (
	"math/rand"
	"time"
)

func main() {
	// 初始化随机数生成器的种子，确保整个程序中的随机数生成器都使用这个种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

}
