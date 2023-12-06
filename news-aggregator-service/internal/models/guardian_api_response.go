package models

type GuardianNewsApiResponse struct {
	Response GuardianResponse `json:"response"`
}

type GuardianResponse struct {
	Status      string           `json:"status"`
	UserTier    string           `json:"userTier"`
	Total       int              `json:"total"`
	StartIndex  int              `json:"startIndex"`
	PageSize    int              `json:"pageSize"`
	CurrentPage int              `json:"currentPage"`
	Pages       int              `json:"pages"`
	Edition     GuardianEdition  `json:"edition"`
	Section     GuardianSection  `json:"section"`
	Results     []GuardianResult `json:"results"`
}

type GuardianEdition struct {
	ID       string `json:"id"`
	WebTitle string `json:"webTitle"`
	WebUrl   string `json:"webUrl"`
	ApiUrl   string `json:"apiUrl"`
	Code     string `json:"code"`
}

type GuardianSection struct {
	ID       string            `json:"id"`
	WebTitle string            `json:"webTitle"`
	WebUrl   string            `json:"webUrl"`
	ApiUrl   string            `json:"apiUrl"`
	Editions []GuardianEdition `json:"editions"`
}

type GuardianResult struct {
	ID                 string         `json:"id"`
	Type               string         `json:"type"`
	SectionId          string         `json:"sectionId"`
	SectionName        string         `json:"sectionName"`
	WebPublicationDate string         `json:"webPublicationDate"`
	WebTitle           string         `json:"webTitle"`
	WebUrl             string         `json:"webUrl"`
	ApiUrl             string         `json:"apiUrl"`
	Fields             GuardianFields `json:"fields"`
	IsHosted           bool           `json:"isHosted"`
	PillarId           string         `json:"pillarId"`
	PillarName         string         `json:"pillarName"`
}

type GuardianFields struct {
	Thumbnail string `json:"thumbnail"`
}
