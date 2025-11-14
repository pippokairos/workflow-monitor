package debug

import "fmt"

var Enabled bool

func Printf(format string, v ...any) {
	if Enabled {
		fmt.Printf("[DEBUG] "+format+"\n", v...)
	}
}
