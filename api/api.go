package api

import (
	"pan-gplimiter/app"

	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

var appSett app.AppSettStruct

type User struct {
	Primary   string
	Domain    string
	Username  string
	Computer  string
	LoginTime int
}

type UserSlice []User

func (a UserSlice) Len() int           { return len(a) }
func (a UserSlice) Less(i, j int) bool { return a[i].LoginTime < a[j].LoginTime }
func (a UserSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (duplicateUsers UserSlice) sortByTime() {

	sort.Sort(duplicateUsers)

}

func RunAPIJobs(appSettParam app.AppSettStruct) {

	appSett = appSettParam

	appSettJSON, err := json.Marshal(appSett)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("RUNNING CONFIG: " + string(appSettJSON))

	duplicateUsers := getDuplicateUsers()

	duplicateUsers.sortByTime()

	duplicateUsersJSON, err := json.Marshal(duplicateUsers)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("DUPPLICATE LIST: " + string(duplicateUsersJSON))

	usersToKick := duplicateUsers.findUsersToKick()

	usersToKickJSON, err := json.Marshal(usersToKick)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("KICK LIST: " + string(usersToKickJSON))

	if !appSett.DryRun {
		usersToKick.kickAll()
	}
}

func (user User) kickUser() {

	var cmd string

	if user.Domain != "" {

		cmd = "<request><global-protect-gateway><client-logout><gateway>" + appSett.GPGateway + "-N</gateway><reason>force-logout</reason><user>" + user.Username + "</user><computer>" + user.Computer + "</computer><domain>" + user.Domain + "</domain></client-logout></global-protect-gateway></request>"

	} else {

		cmd = "<request><global-protect-gateway><client-logout><gateway>" + appSett.GPGateway + "-N</gateway><reason>force-logout</reason><user>" + user.Username + "</user><computer>" + user.Computer + "</computer></client-logout></global-protect-gateway></request>"
	}

	resp := callAPI(cmd)

	xml, err := xmlquery.Parse(strings.NewReader(string(resp)))
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	result := xmlquery.FindOne(xml, "/response[@status=\"success\"]/result/response[@status=\"success\"]")

	if result == nil {

		whiteSpaceRegex := regexp.MustCompile(`\s+`)
		respString := whiteSpaceRegex.ReplaceAllString(string(resp), " ")
		app.LogWarn.Println("KICK ERROR: " + string(respString))

	} else {

		KickedUserJSON, err := json.Marshal(user)
		if err != nil {
			app.LogErr.Fatalln(err)
		}
		app.LogInfo.Println("KICKED LOGIN: " + string(KickedUserJSON))

	}

}

func (usersToKick UserSlice) kickAll() {

	for _, item := range usersToKick {

		item.kickUser()

	}

}

func (duplicateUsers UserSlice) findUsersToKick() (usersToKick UserSlice) {

	primaryMap := make(map[string]int)

	for _, item := range duplicateUsers {

		count, exist := primaryMap[item.Primary]

		if exist {
			primaryMap[item.Primary] += 1
		} else {
			primaryMap[item.Primary] = 1
			continue
		}

		_, found := app.FindString(appSett.ExcludedUsers, item.Primary)

		if found {

			ExcludedUserJSON, err := json.Marshal(item)
			if err != nil {
				app.LogErr.Fatalln(err)
			}
			app.LogInfo.Println("EXCLUDED LOGIN: " + string(ExcludedUserJSON))

			continue
		}

		if count+1 > appSett.MaxLogin {
			usersToKick = append(usersToKick, item)
		}

	}

	return usersToKick
}

func getDuplicateUsers() (duplicateUsers UserSlice) {

	cmd := "<show><global-protect-gateway><current-user><gateway>" + appSett.GPGateway + "</gateway></current-user></global-protect-gateway></show>"
	resp := callAPI(cmd)

	xml, err := xmlquery.Parse(strings.NewReader(string(resp)))
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	result := xmlquery.FindOne(xml, "/response[@status=\"success\"]/result")

	if result == nil {

		whiteSpaceRegex := regexp.MustCompile(`\s+`)
		respString := whiteSpaceRegex.ReplaceAllString(string(resp), " ")
		app.LogErr.Fatalln(string(respString))

	}

	for _, iter := range xmlquery.Find(result, "/entry[primary-username][domain][username][computer][login-time-utc][primary-username=following-sibling::entry/primary-username or primary-username=preceding-sibling::entry/primary-username]") {

		primary := iter.SelectElement("primary-username").InnerText()
		domain := iter.SelectElement("domain").InnerText()
		username := iter.SelectElement("username").InnerText()
		computer := iter.SelectElement("computer").InnerText()
		logintime, err := strconv.Atoi(iter.SelectElement("login-time-utc").InnerText())
		if err != nil {
			app.LogErr.Fatalln(err)
		}

		duplicateUsers = append(duplicateUsers, User{Primary: primary, Domain: domain, Username: username, Computer: computer, LoginTime: logintime})

	}

	return duplicateUsers

}

func callAPI(cmd string) (xml []byte) {

	transport := &(*http.DefaultTransport.(*http.Transport))
	if appSett.SkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{Transport: transport, Timeout: timeout}

	url := "https://" + appSett.FirewallHost + "/api/?key=" + appSett.ApiKey + "&type=op&vsys=vsys" + strconv.Itoa(appSett.VsysNo) + "&cmd=" + url.PathEscape(cmd)

	response, err := client.Get(url)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	return body

}
