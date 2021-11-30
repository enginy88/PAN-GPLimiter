package api

import (
	"pan-gplimiter/app"

	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/xmlquery"
	"github.com/panjf2000/ants"
)

var appSett app.AppSettStruct

type taskFunc func()

type Session struct {
	Primary   string
	Domain    string
	Username  string
	Computer  string
	LoginTime int
}

type SessionSlice []Session

func (a SessionSlice) Len() int           { return len(a) }
func (a SessionSlice) Less(i, j int) bool { return a[i].LoginTime < a[j].LoginTime }
func (a SessionSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (duplicateSessions SessionSlice) sortByTime() {

	sort.Sort(duplicateSessions)

}

func RunAPIJobs(appSettParam app.AppSettStruct) {

	appSett = appSettParam

	appSettJSON, err := json.Marshal(appSett)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("RUNNING CONFIG: " + string(appSettJSON))

	duplicateSessions := getDuplicateSessions()

	duplicateSessions.sortByTime()

	duplicateSessionsJSON, err := json.Marshal(duplicateSessions)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("DUPLICATE LIST: (" + strconv.Itoa(len(duplicateSessions)) + " records) " + string(duplicateSessionsJSON))

	sessionsToKick := duplicateSessions.findSessionsToKick()

	sessionsToKickJSON, err := json.Marshal(sessionsToKick)
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("KICK LIST: (" + strconv.Itoa(len(sessionsToKick)) + " records) " + string(sessionsToKickJSON))

	if !appSett.DryRun {
		sessionsToKick.kickAll()
	}

}

func taskFuncWrapper(session Session, wg *sync.WaitGroup) taskFunc {

	return func() {

		session.kickSession()
		wg.Done()

	}

}

func (session Session) kickSession() {

	var cmd string

	if session.Domain != "" {

		cmd = "<request><global-protect-gateway><client-logout><gateway>" + appSett.GPGateway + "-N</gateway><reason>force-logout</reason><user>" + session.Username + "</user><computer>" + session.Computer + "</computer><domain>" + session.Domain + "</domain></client-logout></global-protect-gateway></request>"

	} else {

		cmd = "<request><global-protect-gateway><client-logout><gateway>" + appSett.GPGateway + "-N</gateway><reason>force-logout</reason><user>" + session.Username + "</user><computer>" + session.Computer + "</computer></client-logout></global-protect-gateway></request>"
	}

	ok, resp := callAPI(cmd, appSett.FailOnError)

	if !ok {
		app.LogWarn.Println("KICK ERROR: " + string(resp))
		return
	}

	xml, err := xmlquery.Parse(strings.NewReader(string(resp)))
	if err != nil {
		if appSett.FailOnError == true {
			app.LogErr.Fatalln(err)
		}
		app.LogWarn.Println("KICK ERROR: " + err.Error() + " RESPONSE: " + string(resp))
		return
	}

	result := xmlquery.FindOne(xml, "/response[@status=\"success\"]/result/response[@status=\"success\"]")

	if result == nil {

		whiteSpaceRegex := regexp.MustCompile(`\s+`)
		respString := whiteSpaceRegex.ReplaceAllString(string(resp), " ")
		app.LogWarn.Println("KICK ERROR: " + string(respString))

	} else {

		KickedSessionJSON, err := json.Marshal(session)
		if err != nil {
			app.LogErr.Fatalln(err)
		}
		app.LogInfo.Println("KICKED LOGIN: " + string(KickedSessionJSON))

	}

}

func (sessionsToKick SessionSlice) kickAll() {

	if appSett.MultiThread {

		var wg sync.WaitGroup
		wp, _ := ants.NewPool(10)
		defer wp.Release()

		for _, session := range sessionsToKick {

			wg.Add(1)
			wp.Submit(taskFuncWrapper(session, &wg))

		}

		wg.Wait()

	} else {

		for _, session := range sessionsToKick {

			session.kickSession()

		}

	}

}

func (duplicateSessions SessionSlice) findSessionsToKick() (sessionsToKick SessionSlice) {

	primaryMap := make(map[string]int)

	for _, session := range duplicateSessions {

		count, exist := primaryMap[session.Primary]

		if exist {
			primaryMap[session.Primary] += 1
		} else {
			primaryMap[session.Primary] = 1
			continue
		}

		_, found := app.FindString(appSett.ExcludedUsers, session.Primary)

		if found {

			ExcludedUsersJSON, err := json.Marshal(session)
			if err != nil {
				app.LogErr.Fatalln(err)
			}
			app.LogInfo.Println("EXCLUDED LOGIN: " + string(ExcludedUsersJSON))

			continue
		}

		if count+1 > appSett.MaxLogin {
			sessionsToKick = append(sessionsToKick, session)
		}

	}

	return sessionsToKick
}

func (allSessions SessionSlice) calculateDuplicates() (duplicateSessions SessionSlice) {

	for i, session := range allSessions {

		found := false

		for _, s := range allSessions[:i] {
			if s.Primary == session.Primary {
				found = true
				break
			}
		}

		if found {
			duplicateSessions = append(duplicateSessions, session)
			continue
		}

		for _, s := range allSessions[i+1:] {
			if s.Primary == session.Primary {
				found = true
				break
			}
		}

		if found {
			duplicateSessions = append(duplicateSessions, session)
			continue
		}

	}

	return duplicateSessions

}

func getDuplicateSessions() (duplicateSessions SessionSlice) {

	var allSessions SessionSlice
	var input io.Reader
	var resp []byte

	if appSett.FileTest == true {
		file, err := os.Open("./test.xml")
		if err != nil {
			app.LogErr.Fatalln(err)
		}
		defer file.Close()
		input = file

	} else {
		cmd := "<show><global-protect-gateway><current-user><gateway>" + appSett.GPGateway + "</gateway></current-user></global-protect-gateway></show>"
		_, resp := callAPI(cmd, true)
		input = strings.NewReader(string(resp))

	}

	xml, err := xmlquery.Parse(input)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	result := xmlquery.FindOne(xml, "/response[@status=\"success\"]/result")

	if result == nil {

		if appSett.FileTest == true {
			app.LogErr.Fatalln("Not an expected XML File!")
		} else {
			whiteSpaceRegex := regexp.MustCompile(`\s+`)
			respString := whiteSpaceRegex.ReplaceAllString(string(resp), " ")
			app.LogErr.Fatalln(respString)
		}

	}

	for _, iter := range xmlquery.Find(result, "/entry[primary-username][domain][username][computer][login-time-utc]") {

		primary := iter.SelectElement("primary-username").InnerText()
		domain := iter.SelectElement("domain").InnerText()
		username := iter.SelectElement("username").InnerText()
		computer := iter.SelectElement("computer").InnerText()
		logintime, err := strconv.Atoi(iter.SelectElement("login-time-utc").InnerText())
		if err != nil {
			app.LogErr.Fatalln(err)
		}

		allSessions = append(allSessions, Session{Primary: primary, Domain: domain, Username: username, Computer: computer, LoginTime: logintime})

	}

	duplicateSessions = allSessions.calculateDuplicates()

	return duplicateSessions

}

func callAPI(cmd string, failOnError bool) (ok bool, xml []byte) {

	transport := &(*http.DefaultTransport.(*http.Transport))
	if appSett.SkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{Transport: transport, Timeout: timeout}

	url := "https://" + appSett.FirewallHost + "/api/?key=" + appSett.ApiKey + "&type=op&vsys=vsys" + strconv.Itoa(appSett.VsysNo) + "&cmd=" + url.PathEscape(cmd)

	response, err := client.Get(url)
	if err != nil {
		if failOnError == true {
			app.LogErr.Fatalln(err)
		}
		return false, []byte(err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	return true, body

}
