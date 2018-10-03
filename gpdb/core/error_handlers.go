package core

import "github.com/op/go-logging"

var (
	log = logging.MustGetLogger("gpdb")
)


// Fatal Handler that terminates the program when its called
func Fatal_handler(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

// Error Handler that return the error when called but doesn't terminate the program
func Error_handler(err error) {
	if err != nil {
		log.Errorf("%s", err)
	}
}

// Warning  Handler that return the warning when called but doesn't terminate the program
func Warn_handler(err error) {
	if err != nil {
		log.Warningf("%s", err)
	}
}