package moodlegrab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GrabClient struct {
	MoodleUrl string
	UserName  string
	Passwd    string
	Client    http.Client

	LoginToken    string
	MoodleSession string
	SessKey       string
}

func (g *GrabClient) ParseSessKey(resp http.Response) error {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Failed to parse HTML:", err)
	}
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()
		if strings.Contains(scriptContent, "sesskey") {
			re := regexp.MustCompile(`"sesskey":"(.*?)"`)
			match := re.FindStringSubmatch(scriptContent)
			if len(match) > 1 {
				g.SessKey = match[1]
			}
		}
	})
	log.Println("sesskey", g.SessKey)
	return nil
}
func (g *GrabClient) makereq(endpoint, method string, data string, isJson bool) (*http.Request, error) {
	req, err := http.NewRequest(method, g.MoodleUrl+endpoint, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Accept-Encoding", "*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	if isJson {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	}
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)

	if g.MoodleSession != "" {
		req.Header.Set("Cookie", "MoodleSession="+g.MoodleSession)
	}
	return req, nil
}
func (g *GrabClient) fetchOriginCookiesAndToken() error {
	req, err := g.makereq("/login/index.php", "GET", "", false)
	if err != nil {
		return err
	}
	resp, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		log.Printf("Cookie: %s = %s", cookie.Name, cookie.Value)
		if cookie.Name == "MoodleSession" {
			g.MoodleSession = cookie.Value
		}
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to parse response body: %v", err)
	}
	token, exists := doc.Find("input[name='logintoken']").Attr("value")
	if !exists {
		return fmt.Errorf("logintoken not exist")
	}
	g.LoginToken = token
	log.Println("logintoken", token)
	return nil
}
func (g *GrabClient) GrepCourses() error {
	data := `[{"index":0,"methodname":"core_course_get_enrolled_courses_by_timeline_classification","args":{"offset":0,"limit":0,"classification":"all","sort":"fullname","customfieldname":"","customfieldvalue":""}}]`
	req, err := g.makereq(fmt.Sprintf("/lib/ajax/service.php?sesskey=%s&info=core_course_get_enrolled_courses_by_timeline_classification", g.SessKey), "POST", data, true)
	resp, err := g.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute POST request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error parsing body")
	}
	var Moodlewp []MoodleJson
	err = json.Unmarshal([]byte(body), &Moodlewp)
	if Moodlewp[0].Error {
		return fmt.Errorf("Error in core_course_get_enrolled_courses_by_timeline_classification. %v", body)
	}
	for _, course := range Moodlewp[0].Data.Courses {
		fmt.Println("-----------------------------------------")
		fmt.Printf("Course ID: %d\n", course.ID)
		fmt.Printf("Full Name: %s\n", course.Fullname)
		fmt.Printf("Short Name: %s\n", course.Shortname)
		fmt.Printf("Start Date: %d\n", course.StartDate)
		fmt.Printf("End Date: %d\n", course.EndDate)
		fmt.Printf("View URL: %s\n", course.ViewURL)
		fmt.Printf("Course Category: %s\n", course.CourseCategory)
		fmt.Println()
		req, err = g.makereq(fmt.Sprintf("/lib/ajax/service.php?sesskey=%s&info=core_courseformat_get_state", g.SessKey), "POST", fmt.Sprintf(`[ { "index": 0, "methodname": "core_courseformat_get_state", "args": { "courseid": %d } } ]`, course.ID), true)
		resp, err = g.Client.Do(req)
		if err != nil {
			return fmt.Errorf("Failed to get core_courseformat_get_state: %v", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error parsing body")
		}
		var Moodlejsonjson []MoodleJsonJson
		err = json.Unmarshal([]byte(body), &Moodlejsonjson)
		g.Parse1Course(Moodlejsonjson[0].Data)
	}

	return nil
}
func (g *GrabClient) Parse1Course(data string) {
	var Moodlejsonjson DetailedMoodleJson
	err := json.Unmarshal([]byte(data), &Moodlejsonjson)
	if err != nil {
		log.Println(err)
	}
	cmMap := make(map[string]DetailedMoodleJson_Cm)
	for _, cm := range Moodlejsonjson.Cm {
		cmMap[cm.ID] = cm
	}
	for _, DetailSection := range Moodlejsonjson.Section {
		fmt.Printf("Title :%s\n", DetailSection.Title)
		for _, cm := range DetailSection.CmList {
			fmt.Printf("        CM :%s\n", cm)
			if cm_, found := cmMap[cm]; found {
				fmt.Printf("        url :%s\n", cm_.URL)
				req, err := g.makereq(strings.Replace(cm_.URL, g.MoodleUrl, "", 1), "GET", "", false)
				resp, err := g.Client.Do(req)
				if err != nil {
					log.Printf("Failed to get URL %s: %v", cm_.URL, err)
					continue
				}
				defer resp.Body.Close()
				finalURL := resp.Request.URL.String()
				if resp.StatusCode == http.StatusSeeOther {
					finalURL = resp.Header.Get("Location")
					fmt.Printf("        final url :%s\n", finalURL)
				}
			}
		}
	}
}
func (g *GrabClient) Login() error {
	err := g.fetchOriginCookiesAndToken()
	if err != nil {
		return fmt.Errorf("Error in Fetch Origin Cookies And Token")
	}
	data := url.Values{}
	data.Set("username", g.UserName)
	data.Set("password", g.Passwd)
	data.Set("logintoken", g.LoginToken)
	data.Set("anchor", "")

	req, err := g.makereq("/login/index.php", "POST", data.Encode(), false)

	resp, err := g.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute POST request: %v", err)
	}
	log.Println("StatCode = ", resp.StatusCode)
	defer resp.Body.Close()
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		log.Printf("Cookie: %s = %s", cookie.Name, cookie.Value)
		if cookie.Name == "MoodleSession" {
			g.MoodleSession = cookie.Value
		}
	}

	req, err = g.makereq("", "GET", "", false)
	resp, err = g.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute GET request: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Login successful")
	err = g.ParseSessKey(*resp)
	return nil
}
