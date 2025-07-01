package app_logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)
var logger *zap.Logger
func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	// Definir la configuración del Logger
	// 0. Salida por consola
	stdout := zapcore.AddSync(os.Stdout)
	// 1. Rotación de ficheros con Lumberjack
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    5,
		MaxBackups: 10,
		MaxAge:     14,
		Compress:   true,
	})
	// 2. Nivel de logging
	level := zap.InfoLevel
	logLevel := zap.NewAtomicLevelAt(level)

	// Production env
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.ConsoleSeparator = ";"
	// Dev env
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder


	// 3. Codificadores
    consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	// var gitRevision string

	// buildInfo, ok := debug.ReadBuildInfo()
	// if ok {
	// 	for _, v := range buildInfo.Settings {
	// 		if v.Key == "vcs.revision" {
	// 			gitRevision = v.Value
	// 			break
	// 		}
	// 	}
	// }
	// 4. Configuración del ZAPCORE
	core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel),
			zapcore.NewCore(fileEncoder, file, logLevel),
				// With(
				// 	[]zapcore.Field{
				// 		zap.String("git_revision", gitRevision),
				// 		zap.String("go_version", buildInfo.GoVersion),
				// 	},
				// )
				
		)
	
	logger = zap.New(core)
	return logger

}


type ctxKey struct{}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
