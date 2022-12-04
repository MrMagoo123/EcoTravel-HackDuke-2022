package funcs

import "fmt"

var (
	LogChan = make(chan Log, 100)
)

type Log struct {
	Color  string
	ID     string
	Msg    string
	Filter string
	Msg2   string
}

func R2(id string, i interface{}) { LogChan <- SingleColorLog("r", id, i) }
func B2(id string, i interface{}) { LogChan <- SingleColorLog("b", id, i) }
func G2(id string, i interface{}) { LogChan <- SingleColorLog("g", id, i) }
func Y2(id string, i interface{}) { LogChan <- SingleColorLog("y", id, i) }
func C2(id string, i interface{}) { LogChan <- SingleColorLog("c", id, i) }
func M2(id string, i interface{}) { LogChan <- SingleColorLog("m", id, i) }
func W2(id string, i interface{}) { LogChan <- SingleColorLog("w", id, i) }

func SingleColorLogMode(mode, color, id string, i interface{}) Log {
	return Log{Color: color, ID: id, Msg: fmt.Sprintf("%v", i)}
}

func SingleColorLog(color, id string, i interface{}) Log {
	if color == "r" {
		return Log{Color: color, ID: id, Msg: fmt.Sprintf("%v", i)}
	}
	return Log{Color: color, ID: id, Msg: fmt.Sprintf("%v", i)}
}

const LOGF = "%v\n"

// --- Standard log (No special formatting)
func R(i interface{})      { LogChan <- Log{Color: "r", Msg: fmt.Sprintf("%v", i), Filter: "ERROR"} }
func B(i interface{})      { LogChan <- Log{Color: "b", Msg: fmt.Sprintf("%v", i)} }
func G(i interface{})      { LogChan <- Log{Color: "g", Msg: fmt.Sprintf("%v", i)} }
func Y(i interface{})      { LogChan <- Log{Color: "y", Msg: fmt.Sprintf("%v", i)} }
func C(i interface{})      { LogChan <- Log{Color: "c", Msg: fmt.Sprintf("%v", i)} }
func M(i interface{})      { LogChan <- Log{Color: "m", Msg: fmt.Sprintf("%v", i)} }
func W(i interface{})      { LogChan <- Log{Color: "w", Msg: fmt.Sprintf("%v", i)} }
func Silent(i interface{}) { LogChan <- Log{Color: "silent", Msg: fmt.Sprintf("%v", i)} }
