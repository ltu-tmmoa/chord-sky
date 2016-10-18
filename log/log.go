package log

import (
	golog "log"
	"os"
)

// Logger is the default application logger.
var Logger *golog.Logger

func init() {
	Logger = golog.New(os.Stdout, "chord-sky ", golog.Ltime | golog.Ldate | golog.Lshortfile | golog.LUTC)
}