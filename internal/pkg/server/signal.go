// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"os"
	"os/signal"
)

// onlyOneSignalHandler 传递信号用，通知 SetupSignalHandler 是否只使用一次（通道的用例！）
var onlyOneSignalHandler = make(chan struct{})

var shutdownHandler chan os.Signal

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.（
func SetupSignalHandler() <-chan struct{} {
	// 保证 iam-apiserver 组件的代码只调用一次 SetupSignalHandler()
	close(onlyOneSignalHandler) // panics when called twice

	// 缓冲区大小为 2
	shutdownHandler = make(chan os.Signal, 2)

	// 使用 struct{} 类型，该类型不占用任何内存空间，作为一种轻量级的信号传递机制（用作信号通知）
	// 其目的仅用来进行信号传递，不传递实际数据
	stop := make(chan struct{})

	// 注册通道，仅用来接收 os.Interrupt 和 syscall.SIGTERM 信号
	signal.Notify(shutdownHandler, shutdownSignals...)

	// 收到一次 SIGINT/ SIGTERM 信号，程序优雅关闭
	// 收到两次 SIGINT/ SIGTERM 信号，程序强制关闭
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	// 返回仅用来发送信号的切片
	return stop
}

// RequestShutdown emulates a received event that is considered as shutdown signal (SIGTERM/SIGINT)
// This returns whether a handler was notified.
func RequestShutdown() bool {
	if shutdownHandler != nil {
		select {
		case shutdownHandler <- shutdownSignals[0]:
			return true
		default:
		}
	}

	return false
}
