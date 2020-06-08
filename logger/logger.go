package logger

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/thetreep/covidtracker"
)

const (
	contextKeyPrefix       = "thetreep.covidtracker."
	ContextKeyRequestId    = contextKeyPrefix + "request_id"
	ContextKeyRequestURI   = contextKeyPrefix + "request_uri"
	ContextKeyRequestLogin = contextKeyPrefix + "request_login"
)

type Logger struct{}

var (
	DefaultLogger = logrus.New()
)

var _ covidtracker.Logfer = &Logger{}

func init() {
	DefaultLogger.Level = logrus.DebugLevel
	//TODO create and set env variable
	if lvl := os.Getenv("THETREEP_COVIDTRACKER_LOG_LEVEL"); lvl != "" {
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

func (l *Logger) HasErr(ctx context.Context, err error) bool {
	return l.HasErrWithFields(ctx, nil, err)
}

func (l *Logger) HasErrWithFields(ctx context.Context, fields map[string]interface{}, err error) bool {
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

func (l *Logger) Debug(ctx context.Context, str string, vars ...interface{}) {
	l.DebugWithFields(ctx, nil, str, vars...)
}

func (l *Logger) DebugWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	l.entry(ctx, fields, false).Debugf(str, vars...)
}

func (l *Logger) Info(ctx context.Context, str string, vars ...interface{}) {
	l.InfoWithFields(ctx, nil, str, vars...)
}

func (l *Logger) InfoWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	l.entry(ctx, fields, false).Infof(str, vars...)
}

func (l *Logger) Warn(ctx context.Context, str string, vars ...interface{}) {
	l.WarnWithFields(ctx, nil, str, vars...)
}

func (l *Logger) WarnWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	l.entry(ctx, fields, false).Warnf(str, vars...)
}

func (l *Logger) Error(ctx context.Context, str string, vars ...interface{}) {
	l.ErrorWithFields(ctx, nil, str, vars...)
}

func (l *Logger) ErrorWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
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

	l.entry(ctx, fields, true).Errorf(str, vars...)
}

func (l *Logger) Panic(ctx context.Context, str string, vars ...interface{}) {
	l.PanicWithFields(ctx, nil, str, vars...)
}

func (l *Logger) PanicWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{}) {
	l.entry(ctx, fields, true).Panicf(str, vars...)
}

func (l *Logger) entry(ctx context.Context, fields map[string]interface{}, includeStack bool) *logrus.Entry {
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
