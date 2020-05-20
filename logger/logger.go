package logger

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

const (
	contextKeyPrefix       = "thetreep.covidtracker."
	ContextKeyRequestId    = contextKeyPrefix + "request_id"
	ContextKeyRequestURI   = contextKeyPrefix + "request_uri"
	ContextKeyRequestLogin = contextKeyPrefix + "request_login"
)

var (
	DefaultLogger = logrus.New()
)

func init() {
	DefaultLogger.Level = logrus.DebugLevel
	//TODO create and set env variable
	if lvl := os.Getenv("THETREEP_COVID_TRACKER_LOG_LEVEL"); lvl != "" {
		parsed, err := logrus.ParseLevel(lvl)
		if err != nil {
			fmt.Printf("invalid loging level %q\n", lvl)
		} else {
			DefaultLogger.Level = parsed
		}
	}
	DefaultLogger.Out = os.Stdout
	DefaultLogger.Formatter = &logrus.JSONFormatter{}
}

func HasErr(ctx context.Context, err error) bool {
	return HasErrWithFields(ctx, nil, err)
}

func HasErrWithFields(ctx context.Context, fields map[string]interface{}, err error) bool {
	if err != nil {
		if fields == nil {
			fields = make(map[string]interface{})
		}
		if tags, ok := fields["tags"]; ok {
			switch tt := tags.(type) {
			case []interface{}:
				fields["tags"] = append(tt, "error")
			default:
				fields["tags"] = []interface{}{tt, "error"}
			}
		} else {
			fields["tags"] = "error"
		}
		entry(ctx, fields, true).Warnf("error: %s", err) // Normally, all errors here should be at least with status Warn
		return true
	}
	return false
}

func Debug(ctx context.Context, str string, vars ...interface{}) {
	DebugWithFields(ctx, nil, str, vars...)
}

func DebugWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	entry(ctx, fields, false).Debugf(str, vars...)
}

func Info(ctx context.Context, str string, vars ...interface{}) {
	InfoWithFields(ctx, nil, str, vars...)
}

func InfoWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	entry(ctx, fields, false).Infof(str, vars...)
}

func Warn(ctx context.Context, str string, vars ...interface{}) {
	WarnWithFields(ctx, nil, str, vars...)
}

func WarnWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	entry(ctx, fields, false).Warnf(str, vars...)
}

func Error(ctx context.Context, str string, vars ...interface{}) {
	ErrorWithFields(ctx, nil, str, vars...)
}

func ErrorWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if tags, ok := fields["tags"]; ok {
		switch tt := tags.(type) {
		case []interface{}:
			fields["tags"] = append(tt, "error")
		default:
			fields["tags"] = []interface{}{tt, "error"}
		}
	} else {
		fields["tags"] = "error"
	}

	entry(ctx, fields, true).Errorf(str, vars...)
}

func Panic(ctx context.Context, str string, vars ...interface{}) {
	PanicWithFields(ctx, nil, str, vars...)
}

func PanicWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	entry(ctx, fields, true).Panicf(str, vars...)
}

func entry(ctx context.Context, fields map[string]interface{}, includeStack bool) *logrus.Entry {
	entry := DefaultLogger.WithFields(fields)
	if id := ctx.Value(ContextKeyRequestId); id != nil {
		entry = entry.WithField("RequestID", id)
	}
	if uri := ctx.Value(ContextKeyRequestURI); uri != nil {
		entry = entry.WithField("RequestURI", uri)
	}
	if login := ctx.Value(ContextKeyRequestLogin); login != nil {
		entry = entry.WithField("RequestLogin", login)
	}
	if includeStack {
		entry = entry.WithField("Stack", string(debug.Stack()))
	}
	return entry
}
