package main

import (
	"PAN-GPLimiter/api"
	"PAN-GPLimiter/app"

	"fmt"
	"time"
)

var appFlag app.AppFlagStruct
var appSett app.AppSettStruct

func main() {

	start := time.Now()
	app.LogAlways.Println("HELLO MSG: Welcome to PAN-GPLimiter v2.1 by EY!")

	appFlag = app.GetAppFlag()
	appSett = app.GetAppSett()

	api.RunAPIJobs(appSett)

	duration := fmt.Sprintf("%.1f", time.Since(start).Seconds())
	app.LogAlways.Println("BYE MSG: All done in " + duration + "s, bye!")

}
