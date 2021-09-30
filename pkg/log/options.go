package log

import (
	"io"
	"os"
)

const (
	FmtEmptySeparate = ""
)

type Level uint8

const (
	// DebugLevel 主要用来提供一些Debug信息，方便开发测试时，定位问题，一般量很大
	DebugLevel Level = iota
	// InfoLevel 信息，默认日志等级，提供一些必要的日志信息，方便业务出问题时，结合Error排查故障
	InfoLevel
	// WarnLevel 警告等级比InfoLevel高但一般不需要人工处理
	WarnLevel
	// ErrorLevel 高优先级的错误等级 实际的程序执行中不应该产生此类错误
	ErrorLevel
	// PanicLevel 恐慌错误并调用panic
	PanicLevel
	// FatalLevel 记录致命错误并调用os.Exit(1)
	FatalLevel
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

type options struct {
	output        io.Writer
	level         Level
	stdLevel      Level
	formatter     Formatter
	disableCaller bool
}

type Option func(*options)

func initOptions(opts ...Option) *options {
	o := &options{
		output:    os.Stdout,
		level:     InfoLevel,
		formatter: &TextFormatter{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithStdLevel
func WithStdLevel(level Level) Option {
	return func(o *options) {
		o.stdLevel = level
	}
}

// WithLevel 设置输出级别
func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

// WithOutput 设置输出位置
func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

// WithFormatter 设置输出格式
func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

// WithDisableCaller 设置是否打印文件名和行号
func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller
	}
}
