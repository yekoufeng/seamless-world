package entity

import (
	"fmt"
)

// String Entity基础信息
func (e *Entity) String() string {
	return fmt.Sprintf("[Type:%s EntityID:%d DBID:%d CellID:%d]",
		e.entityType, e.entityID, e.dbid, e.cellID)
}

// Info 日志
func (e *Entity) Info(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Info(params...)
}

// Infof 日志
func (e *Entity) Infof(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Infof(ff, params...)
}

// Warn 日志
func (e *Entity) Warn(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Warn(params...)
}

// Warnf 日志
func (e *Entity) Warnf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Warnf(ff, params...)
}

// Error 日志
func (e *Entity) Error(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Error(params...)
}

func (e *Entity) JcmiaoTempLog(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Error(params...)
}

// Errorf 日志
func (e *Entity) Errorf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Errorf(ff, params...)
}

// Debug 日志
func (e *Entity) Debug(v ...interface{}) {
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Debug(params...)
}

// Debugf 日志
func (e *Entity) Debugf(format string, v ...interface{}) {
	ff := "%s " + format
	params := []interface{}{e.String()}
	params = append(params, v...)
	elogger.Debugf(ff, params...)
}
