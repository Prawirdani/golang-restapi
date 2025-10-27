package log

import (
	"os"
	"testing"

	stdlog "log"

	"github.com/prawirdani/golang-restapi/config"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.LoadConfig("../../.env")
	if err != nil {
		stdlog.Fatal("failed to load config:", err)
	}
	os.Exit(m.Run())
}

// BenchmarkZerolog_Info-8            20172             59824 ns/op             304 B/op          4 allocs/op
// func BenchmarkZerolog_Info(b *testing.B) {
// 	b.ReportAllocs()
//
// 	logger := NewZerologAdapter(cfg)
//
// 	for b.Loop() {
// 		logger.Info("benchmark message", "foo", "bar")
// 	}
// }

// 19593             57208 ns/op             384 B/op          8 allocs/op
// func BenchmarkSlog_Info(b *testing.B) {
// 	b.ReportAllocs()
// 	b.ResetTimer()
// 	logger := NewSlogAdapter(cfg)
//
// 	for i := 0; i < b.N; i++ {
// 		logger.Info("benchmark message", "foo", "bar")
// 	}
// }
