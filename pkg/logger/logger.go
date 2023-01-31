package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultFilePerm = 0o766

	DefaultFileRotateMaxSize    = 30
	DefaultFileRotateMaxBackups = 6
	DefaultFileRotateMaxAge     = 30

	DefaultTimeLayout = time.RFC3339
)

// Option 自定义日志配置
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

// WithDebugLevel 仅输出高于`debug`级别的日志
func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.DebugLevel
	}
}

// WithInfoLevel 仅输出高于`info`级别的日志
func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.InfoLevel
	}
}

// WithWarnLevel 仅输出高于`warn`级别的日志
func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

// WithErrorLevel 仅输出高于`error`级别的日志
func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

// WithField 为日志输出添加一些特定字段
func WithField(key, value string) Option {
	return func(opt *option) {
		opt.fields[key] = value
	}
}

// WithFile 将日志输入至指定文件
func WithFile(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, DefaultFilePerm); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, DefaultFilePerm)
	if err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = zapcore.Lock(f)
	}
}

func GenerateDefaultLBJWriter(file string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,                        // 日志文件路径
		MaxSize:    DefaultFileRotateMaxSize,    // 单个日志文件的最大尺寸，默认单位为M
		MaxBackups: DefaultFileRotateMaxBackups, // 最多保留多少份日志文件
		MaxAge:     DefaultFileRotateMaxAge,     // 日志文件的最大时间跨度，默认单位为day
		LocalTime:  true,                        // 使用本地时间
		Compress:   false,                       // 是否压缩，默认为禁用压缩
	}
}

func GenerateLBJWriter(file string, maxSize, maxBackups, maxAge int) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
		Compress:   false,
	}
}

// WithFileRotation 将日志输入至指定文件(滚动)
func WithFileRotation(file *lumberjack.Logger) Option {
	dir := filepath.Dir(file.Filename)
	if err := os.MkdirAll(dir, DefaultFilePerm); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = file
	}
}

// WithDisableConsole 禁止输出至stdout
func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

// WithTimeLayout 设定日志时间格式
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// NewLogger 返回一个console-encoder风格的zap logger
func NewLogger(opts ...Option) *zap.SugaredLogger {
	opt := &option{level: zapcore.DebugLevel, fields: make(map[string]string)}
	for _, setup := range opts {
		setup(opt)
	}

	if opt.timeLayout == "" {
		opt.timeLayout = DefaultTimeLayout
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(opt.timeLayout))
	}
	consoleEnc := zapcore.NewConsoleEncoder(config)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEnc,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(consoleEnc,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(consoleEnc,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core, zap.AddCaller(), zap.ErrorOutput(stderr))
	for k, v := range opt.fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{Key: k, Type: zapcore.StringerType, String: v}))
	}

	return logger.Sugar()
}
