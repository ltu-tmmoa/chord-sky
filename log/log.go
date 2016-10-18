package log

import (
	"log"
	"os"
)

// Logger is the default application logger.
var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "chord-sky ", log.Ltime | log.Ldate | log.Lshortfile | log.LUTC)
}