package log

type Formatter interface {
	// 可能在异步 goroutine 中
	// 请将结果写入缓冲区
	Format(entry *Entry) error
}
