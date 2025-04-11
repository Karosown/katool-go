package test

import (
	"os"
	"path"
	"testing"

	"github.com/karosown/katool-go/log"
	rm "github.com/karosown/katool-go/net/format"
	remote "github.com/karosown/katool-go/net/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newLogger(c *Config) *zap.SugaredLogger {
	writeSyncer := getLogWriter(c)
	encoder := getEncoder()
	fileCore := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)
	core := zapcore.NewTee(fileCore, consoleCore)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

type Config struct {
	FileDir    string `yaml:"FileDir"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
	Compress   bool   `yaml:"Compress"`
	LocalTime  bool   `yaml:"LocalTime"`
	Level      zapcore.Level
}

func getLogWriter(c *Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(c.FileDir, "current.log"),
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func TestOAuth2WithZapLog(t *testing.T) {
	req := remote.Req{}
	config := Config{
		FileDir:    "log",
		MaxSize:    1,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   false,
	}
	log := newLogger(&config)
	req.SetLogger(log)
	req.Url("http://www.baidu.com").QueryParam(map[string]string{"a": "1"}).Build(&map[string]string{})
}
func TestOAuth2WithLogrus(t *testing.T) {
	req := remote.Req{}
	req.SetLogger(&log.LogrusAdapter{})
	req.Url("http://www.baidu.com").QueryParam(map[string]string{"a": "1"}).Build(&map[string]string{})

	req.Url("http://www.baidu.com").Format(&rm.JSONEnDeCodeFormat{}).Build(&map[string]string{})
}
