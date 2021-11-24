package app

import (
	"fmt"
	"log"
	"os"
)

var LogErr *log.Logger
var LogWarn *log.Logger
var LogInfo *log.Logger
var LogAlways *log.Logger

func init() {

	LogErr = log.New(os.Stderr, "(PAN-GPLIMITER) ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogWarn = log.New(os.Stdout, "(PAN-GPLIMITER) WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogInfo = log.New(os.Stdout, "(PAN-GPLIMITER) INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogAlways = log.New(os.Stdout, "(PAN-GPLIMITER) ALWAYS: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

func Typeof(v interface{}) string {

	return fmt.Sprintf("%T", v)

}

func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func isValidHost(s string) bool {

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '.') || (r == ':')) {
			return false
		}
	}

	return true

}

func isValidBase64(s string) bool {

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '+') || (r == '/') || (r == '=')) {
			return false
		}
	}

	return true

}

func isValidName(s string) bool {

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '-') || (r == '_') || (r == ' ')) {
			return false
		}
	}

	return true

}

func isValidUser(s string) bool {

	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '\\') || (r == '.')) {
			return false
		}
	}

	return true

}
