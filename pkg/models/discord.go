package models

type Embed struct {
	Author struct {
		Name    string `json:"name"`
		IconUrl string `json:"icon_url"`
		Url     string `json:"url"`
	} `json:"author"`
	Footer struct {
		Text    string `json:"text"`
		IconUrl string `json:"icon_url"`
	} `json:"footer"`
	Image struct {
		Url string `json:"url"`
	} `json:"image"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Title       string `json:"title"`
	Thumbnail   struct {
		Url string `json:"url"`
	} `json:"thumbnail"`
	Url    string `json:"url"`
	Fields []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`
	Color int `json:"color"`
}

type Embeds struct {
	Embeds []Embed `json:"embeds"`
}