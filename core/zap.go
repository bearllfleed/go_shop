package core

import (
	"github.com/bearllflee/go_shop/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

// 一个核心只能写一个日志级别
func NewZapCore(level zapcore.Level) *ZapCore {
	entity := &ZapCore{level: level}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer() zapcore.WriteSyncer {
	cutter := NewCutter(
		CutterWithLayout(global.CONFIG.Logger.Layout),
		CutterWithLevel(z.level),
		CutterWithDirector(global.CONFIG.Logger.Direcotr),
	)
	return zapcore.AddSync(cutter)
}
