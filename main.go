package main

import (
	"pan-gplimiter/api"
	"pan-gplimiter/app"

	"fmt"
	"time"
)

var appSett app.AppSettStruct

func main() {

	start := time.Now()
	app.LogInfo.Println("HELLO MSG: Welcome to PAN-GPLimiter v1.5 by EY")

	appSett = app.GetAppSett()

	api.RunAPIJobs(appSett)

	duration := fmt.Sprintf("%.1f", time.Since(start).Seconds())
	app.LogInfo.Println("BYE MSG: All done in " + duration + "s, bye!")

}
