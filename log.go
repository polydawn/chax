package main

import (
	"io"

	"polydawn.net/chax/chaxui"

	"github.com/google/gxui"
	"github.com/inconshreveable/log15"
)

var _ io.Writer = LogWriter{}

type LogWriter struct {
	Log   log15.Logger
	Level log15.Lvl
	Msg   string
}

func (lw LogWriter) Write(msg []byte) (int, error) {
	// goddamnit, expose your fucking `write` method
	switch lw.Level {
	case log15.LvlDebug:
		lw.Log.Debug(lw.Msg, "chunk", string(msg))
	case log15.LvlInfo:
		lw.Log.Info(lw.Msg, "chunk", string(msg))
	case log15.LvlWarn:
		lw.Log.Warn(lw.Msg, "chunk", string(msg))
	case log15.LvlError:
		lw.Log.Error(lw.Msg, "chunk", string(msg))
	case log15.LvlCrit:
		lw.Log.Crit(lw.Msg, "chunk", string(msg))
	}
	return len(msg), nil
}

func ListDumpHandler(driver gxui.Driver, listCtrl *chaxui.LinesAdapter, fmtr log15.Format) log15.Handler {
	return log15.FuncHandler(func(r *log15.Record) error {
		driver.Call(func() {
			listCtrl.Append(string(fmtr.Format(r)))
		})
		return nil
	})
}
