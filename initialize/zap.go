package initialize

import (
	"fmt"

	ccore "github.com/bearllflee/go_shop/core"
	"github.com/bearllflee/go_shop/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func MustLoadZap() {
	levels := Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)

	for i := 0; i < length; i++ {
		core := ccore.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...))
	global.Logger = logger
}

// Levels 根据字符串转化为 zapcore.Levels
func Levels() []zapcore.Level {
	levels := make([]zapcore.Level, 0, 7)
	level, err := zapcore.ParseLevel(global.CONFIG.Logger.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	for ; level <= zapcore.FatalLevel; level++ {
		levels = append(levels, level)
	}
	fmt.Println(levels)
	// 当你的系统设定了一个日志级别后，系统只会生效输出大于等于该日志级别的日志
	// info  warn error fatal panic  // 不会记录的是 debug
	return levels
}
