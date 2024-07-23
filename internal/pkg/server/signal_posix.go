// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package server

import (
	"os"
	"syscall"
)

// 定义终止信号类型的切片
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
