package log

import (
	"io"

	"github.com/inconshreveable/log15"
)

var _ io.Writer = Writer{}

/*
	Exposes a logger as an `io.Writer`, so you can have any system output tagged logs.
*/
type Writer struct {
	Log   log15.Logger
	Level log15.Lvl
	Msg   string
}

func (lw Writer) Write(msg []byte) (int, error) {
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

type Printer func(a ...interface{}) (n int, err error)

/*
	Dumps logs to a printer.  Any regular `Printer` will do -- `fmt.Print`
	qualifies.
*/
func PrinterHandler(printer Printer, fmtr log15.Format) log15.Handler {
	return log15.FuncHandler(func(r *log15.Record) error {
		printer(string(fmtr.Format(r)))
		return nil
	})
}
