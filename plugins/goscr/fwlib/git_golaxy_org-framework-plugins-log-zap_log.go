// Code generated by 'yaegi extract git.golaxy.org/framework/plugins/log/zap_log'. DO NOT EDIT.

package fwlib

import (
	"git.golaxy.org/framework/plugins/log/zap_log"
	"reflect"
)

func init() {
	Symbols["git.golaxy.org/framework/plugins/log/zap_log/zap_log"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Install":             reflect.ValueOf(&zap_log.Install).Elem(),
		"NewConsoleZapLogger": reflect.ValueOf(zap_log.NewConsoleZapLogger),
		"NewJsonZapLogger":    reflect.ValueOf(zap_log.NewJsonZapLogger),
		"Uninstall":           reflect.ValueOf(&zap_log.Uninstall).Elem(),
		"With":                reflect.ValueOf(&zap_log.With).Elem(),

		// type definitions
		"LoggerOptions": reflect.ValueOf((*zap_log.LoggerOptions)(nil)),
	}
}