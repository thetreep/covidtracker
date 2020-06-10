package covidtracker

import "context"

type Logfer interface {
	HasErr(ctx context.Context, err error) bool
	HasErrWithFields(ctx context.Context, fields map[string]interface{}, err error) bool
	Debug(ctx context.Context, str string, vars ...interface{})
	DebugWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Info(ctx context.Context, str string, vars ...interface{})
	InfoWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Warn(ctx context.Context, str string, vars ...interface{})
	WarnWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Error(ctx context.Context, str string, vars ...interface{})
	ErrorWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Panic(ctx context.Context, str string, vars ...interface{})
	PanicWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
}
