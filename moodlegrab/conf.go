package moodlegrab


type MoodleCourse struct {
	ID                    int     `json:"id"`
	Fullname              string  `json:"fullname"`
	Shortname             string  `json:"shortname"`
	IDNumber              string  `json:"idnumber"`
	Summary               string  `json:"summary"`
	SummaryFormat         int     `json:"summaryformat"`
	StartDate             int64   `json:"startdate"`
	EndDate               int64   `json:"enddate"`
	Visible               bool    `json:"visible"`
	ShowActivityDates     bool    `json:"showactivitydates"`
	ShowCompletionConditions interface{} `json:"showcompletionconditions"`
	FullnameDisplay       string  `json:"fullnamedisplay"`
	ViewURL               string  `json:"viewurl"`
	CourseImage           string  `json:"courseimage"`
	Progress              int     `json:"progress"`
	HasProgress           bool    `json:"hasprogress"`
	IsFavourite           bool    `json:"isfavourite"`
	Hidden                bool    `json:"hidden"`
	ShowShortname         bool    `json:"showshortname"`
	CourseCategory        string  `json:"coursecategory"`
}

type MoodleData struct {
	Courses []MoodleCourse `json:"courses"`
}

type MoodleJson struct {
	Error bool       `json:"error"`
	Data  MoodleData `json:"data"`
}
type MoodleJsonJson struct{
	Error bool       `json:"error"`
	Data  string `json:"data"`
}
type DetailedMoodleJson_course struct{
	ID          string      `json:"id"`
	NumSections int      `json:"numsections"`
	SectionList []string `json:"sectionlist"`
	EditMode    bool     `json:"editmode"`
	Highlighted string   `json:"highlighted"`
	MaxSections string      `json:"maxsections"`
	BaseURL     string   `json:"baseurl"`
	StateKey    string   `json:"statekey"`
}
type DetailedMoodleJson_section struct{
	ID              string      `json:"id"`
	Section         int      `json:"section"`
	Number          int      `json:"number"`
	Title           string   `json:"title"`
	HasSummary      bool     `json:"hassummary"`
	RawTitle        string  `json:"rawtitle"`
	CmList          []string `json:"cmlist"`
	Visible         bool     `json:"visible"`
	SectionURL      string   `json:"sectionurl"`
	Current         bool     `json:"current"`
	IndexCollapsed  bool     `json:"indexcollapsed"`
	ContentCollapsed bool    `json:"contentcollapsed"`
	HasRestrictions bool     `json:"hasrestrictions"`
}
type DetailedMoodleJson_Cm struct{
	ID              string    `json:"id"`
	Anchor          string `json:"anchor"`
	Name            string `json:"name"`
	Visible         bool   `json:"visible"`
	Stealth         bool   `json:"stealth"`
	SectionID       string `json:"sectionid"`
	SectionNumber   int    `json:"sectionnumber"`
	UserVisible     bool   `json:"uservisible"`
	HasCmRestrictions bool `json:"hascmrestrictions"`
	Module          string `json:"module"`
	Plugin          string `json:"plugin"`
	Indent          int    `json:"indent"`
	AccessVisible   bool   `json:"accessvisible"`
	URL             string `json:"url"`
	IsTrackedUser   bool   `json:"istrackeduser"`
}
type DetailedMoodleJson struct{
	Course DetailedMoodleJson_course `json:"course"`
	Section []DetailedMoodleJson_section `json:"section"`
	Cm []DetailedMoodleJson_Cm `json:"cm"`
}

type Yaml_config struct {
	General General
}

type General struct {
	Account   string
	Passwd   string
}
