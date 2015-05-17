package chaxui

import (
	"github.com/google/gxui"
	"github.com/inconshreveable/log15"
)

func ListDumpHandler(driver gxui.Driver, listCtrl *LinesAdapter, fmtr log15.Format) log15.Handler {
	return log15.FuncHandler(func(r *log15.Record) error {
		driver.Call(func() {
			listCtrl.Append(string(fmtr.Format(r)))
		})
		return nil
	})
}
