package app

import (
	"io/ioutil"
	"strconv"

	"github.com/JeremyLoy/config"
)

const (
	defaultMaxLogin    = 1
	defaultVsysNo      = 1
	defaultSkipVerify  = false
	defaultLogSilence  = false
	defaultDryRun      = false
	defaultMultiThread = false
)

type AppSettStruct struct {
	FirewallHost  string   `config:"PANGPLIMITER_FIREWALL_HOST"`
	ApiKey        string   `config:"PANGPLIMITER_API_KEY"`
	GPGateway     string   `config:"PANGPLIMITER_GP_GATEWAY"`
	VsysNo        int      `config:"PANGPLIMITER_VSYS_NO"`
	MaxLogin      int      `config:"PANGPLIMITER_MAX_LOGIN"`
	ExcludedUsers []string `config:"PANGPLIMITER_EXCLUDED_USERS"`
	SkipVerify    bool     `config:"PANGPLIMITER_SKIP_VERIFY"`
	LogSilence    bool     `config:"PANGPLIMITER_LOG_SILENCE"`
	DryRun        bool     `config:"PANGPLIMITER_DRY_RUN"`
	MultiThread   bool     `config:"PANGPLIMITER_MULTI_THREAD"`
}

var appSett AppSettStruct

func GetAppSett() AppSettStruct {

	loadAppSett()
	checkAppSett()

	return appSett

}

func loadAppSett() {

	err := config.From("appsett.env").FromEnv().To(&appSett)
	if err != nil {
		LogErr.Fatalln("Cannot find/load 'appsett.env' file! (" + err.Error() + ")")
	}

}

func checkAppSett() {

	if appSett.LogSilence == true {
		LogWarn.SetOutput(ioutil.Discard)
		LogInfo.SetOutput(ioutil.Discard)
	}

	if appSett.DryRun == true {
		LogWarn.Println("CONFIG MSG: DryRun object value set to ('" + strconv.FormatBool(appSett.DryRun) + "'), no users going to be kicked.")
	}

	if appSett.SkipVerify == true {
		LogWarn.Println("CONFIG MSG: SkipVerify object value set to ('" + strconv.FormatBool(appSett.SkipVerify) + "'), which may be led to security risks.")
	}

	if appSett.MultiThread == true {
		LogWarn.Println("CONFIG MSG: MultiThread object value set to ('" + strconv.FormatBool(appSett.MultiThread) + "'), parallel processing is enabled.")
	}

	if len(appSett.FirewallHost) < 1 || len(appSett.FirewallHost) > 31 {
		LogErr.Fatalln("FirewallHost object ('" + appSett.FirewallHost + "') should has more than 1 chars and less than 31 chars!")
	}
	if !isValidHost(appSett.FirewallHost) {
		LogErr.Fatalln("FirewallHost object ('" + appSett.FirewallHost + "') should only contains alphanumeric chars with space & dash & underscore!")
	}

	if len(appSett.ApiKey) < 0 {
		LogErr.Fatalln("ApiKey object ('" + appSett.ApiKey + "') should has more than 1 chars chars!")
	}
	if !isValidBase64(appSett.ApiKey) {
		LogErr.Fatalln("ApiKey object ('" + appSett.ApiKey + "') should only contains alphanumeric chars with plus & slash & equal!")
	}

	if len(appSett.GPGateway) < 0 {
		LogErr.Fatalln("GPGateway object ('" + appSett.GPGateway + "') should has more than 1 chars chars!")
	}
	if !isValidName(appSett.GPGateway) {
		LogErr.Fatalln("GPGateway object ('" + appSett.GPGateway + "') should only contains alphanumeric chars with dot & colon!")
	}

	if appSett.VsysNo != 0 {
		if appSett.VsysNo < 0 || appSett.VsysNo > 255 {
			LogErr.Fatalln("VsysNo object value ('" + strconv.Itoa(appSett.VsysNo) + "') should be between 1 & 255!")
		}
	} else {
		appSett.VsysNo = defaultVsysNo
		LogWarn.Println("CONFIG MSG: Using default value ('" + strconv.Itoa(appSett.VsysNo) + "') for VsysNo object value.")
	}

	if appSett.MaxLogin != 0 {
		if appSett.MaxLogin < 0 || appSett.MaxLogin > 8 {
			LogErr.Fatalln("MaxLogin object value ('" + strconv.Itoa(appSett.MaxLogin) + "') should be between 1 & 8!")
		}
	} else {
		appSett.MaxLogin = defaultMaxLogin
		LogWarn.Println("CONFIG MSG: Using default value ('" + strconv.Itoa(appSett.MaxLogin) + "') for MaxLogin object value.")
	}

	if appSett.ExcludedUsers == nil {
	} else {
		for _, user := range appSett.ExcludedUsers {
			if !isValidUser(user) {
				LogErr.Fatalln("ExcludedUsers object contents ('" + user + "') should only contains alphanumeric chars with backslash & dot!")
			}
		}
		LogInfo.Println("CONFIG MSG: " + strconv.Itoa(len(appSett.ExcludedUsers)) + " user(s) excluded from concurrent login restriction.")
	}

}
