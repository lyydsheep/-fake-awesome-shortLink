package example

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type ZapSuite struct {
	suite.Suite
}

// TestZap 是使用zap的demo
func (s *ZapSuite) TestExample() {
	l := zap.NewExample()
	l.Debug("this is debug message")
	l.Info("this is info message")
	l.Info("this is info message with field",
		zap.Int("age", 1), zap.String("name", "aba"))
	l.Warn("this is warn message")
	l.Error("this is error message")
}

func (s *ZapSuite) TestDevelopment() {
	l, _ := zap.NewDevelopment()
	defer l.Sync()
	l.Info("failed to fetch url",
		zap.String("url", "www.baidu.com"),
		zap.Int("attempt", 3),
		zap.Duration("duration", time.Second))

	l.With(zap.String("url", "www.baidu.com"),
		zap.Int("attempt", 4),
		zap.Duration("duration", time.Second*5)).Info("[With] failed to fetch url")
}

func Test(t *testing.T) {
	suite.Run(t, new(ZapSuite))
}
