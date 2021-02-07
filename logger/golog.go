package loger

import (
	"github.com/curltech/go-colla-core/config"
	"github.com/curltech/go-colla-core/logger"
)

func init() {
	level, _ := config.GetString("log.level", "debug")
	logger.SetLevel(level)
	timeFormat := config.AppParams.TimeFormat
	logger.SetTimeFormat(timeFormat)
	// Levels contains a map of the log levels and their attributes.
	errorAttrs := logger.Levels[logger.Sugar.ErrorLevel]

	// Change a log level's text.
	customColorCode := 156
	errorAttrs.SetText("custom text", customColorCode)

	// Get (rich) text per log level.
	enableColors := true
	errorAttrs.Text(enableColors)
	logger.Sugar.Infof("log config completed!")
}
