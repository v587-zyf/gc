package utils

import (
	"fmt"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

// WaitExit 等待退出(cirt+c有效,一定加在main函数的最后一行)
func WaitExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("press cirt+c exit app")
	select {
	case sig := <-c:
		fmt.Println(sig, "exit app")
	}
}

var bufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 4096)
		return buf
	},
}

func GoSafe(fnName string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			buf := bufferPool.Get().(*[]byte)
			l := runtime.Stack(*buf, false)
			stackTrace := strings.TrimSuffix(string(string(*buf)[:l]), "\n")
			logMsg := fmt.Sprintf("Recovered from panic in %s", fnName)
			log.Error(logMsg,
				zap.String("panic", fmt.Sprint(r)),
				zap.String("stack", stackTrace),
			)
			bufferPool.Put(buf)
		}
	}()

	go fn()
}
