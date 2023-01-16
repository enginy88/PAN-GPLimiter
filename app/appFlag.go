package app

import (
	"flag"
	"os"
)

type AppFlagStruct struct {
	WorkingDir string
}

var appFlag AppFlagStruct

func GetAppFlag() AppFlagStruct {

	parseAppFlag()
	changeWorkingDir()

	return appFlag

}

func parseAppFlag() {

	workingDir := flag.String("dir", "", "Path of directory which contains 'appsett.env' file.")
	flag.Parse()

	appFlag.WorkingDir = *workingDir

}

func changeWorkingDir() {

	origDir, err := os.Getwd()
	if err != nil {
		LogErr.Fatalln("Cannot get working directory! (" + err.Error() + ")")
	}

	if appFlag.WorkingDir != "" {

		err := os.Chdir(appFlag.WorkingDir)
		if err != nil {
			LogErr.Fatalln("Cannot change working directory! (" + err.Error() + ")")
		}

		newDir, err := os.Getwd()
		if err != nil {
			LogErr.Fatalln("Cannot get working directory! (" + err.Error() + ")")
		}

		LogInfo.Println("CONFIG MSG: Flag 'dir' set, changing working directory from '" + origDir + "' to '" + newDir + "'.")
	}

}
