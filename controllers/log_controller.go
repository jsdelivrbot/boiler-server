package controllers

import (
	"github.com/AzureRelease/boiler-server/models/logs"
	"github.com/AzureRelease/boiler-server/conf"
	"errors"
)

type LogController struct {
	MainController
}

var LogCtl *LogController = &LogController{}

func (ctl *LogController) AddReloadLog(lg *logs.BoilerRuntimeLog) error {
	if !conf.IsReloadLogEnabled {
		return errors.New("reload log is disabled")
	}

	var err error

	go func() {
		err = DataCtl.AddData(lg, false)
	}()

	return err
}