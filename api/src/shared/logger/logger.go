package logger

import "go.uber.org/zap"

var L *zap.SugaredLogger

func init() {
	L = zap.Must(zap.NewDevelopment()).Sugar()
}
