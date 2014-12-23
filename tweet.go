package main

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id        int64  `json:"id"`
	IdStr     string `json:"id_str"`
	Text      string `json:"text"`
	User      User   `json:"user"`
}

type Entities struct {
	Hashtags     []Hashtag     `json:"hashtags"`
	Media        []Media       `json:"media"`
	URLs         []URL         `json:"urls"`
	UserMentions []UserMention `json:"user_mentions"`
}

type Hashtag struct {
	Indices []int  `json:"indicies"`
	Text    string `json:"string"`
}

type Media struct {
	DisplayURL        string                 `json:"display_url"`
	ExpandedURL       string                 `json:"expanded_url"`
	Id                int64                  `json:"id"`
	IdStr             string                 `json:"id_str"`
	Indices           []int                  `json:"indices"`
	MediaURL          string                 `json:"media_url"`
	MediaURLHTTPS     string                 `json:"media_url_https"`
	Sizes             map[string]interface{} `json:"sizes"`
	SourceStatusId    int64                  `json:"source_status_id"`
	SourceStatudIdStr string                 `json:"source_status_id_str"`
	Type              string                 `json:"type"`
	URL               string                 `json:"url"`
}

type URL struct {
	DisplayURL  string `json:"display_url"`
	ExpandedURL string `json:"expanded_url"`
	Indices     []int  `json:"indices"`
	URL         string `json:"url"`
}

type UserMention struct {
	Id         int64  `json:"id"`
	IdStr      string `json:"id_str"`
	Indices    []int  `json:"indices"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type User struct {
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

type EventType string

const (
	AccessRevoked EventType = "access_revoked"
	Favorite      EventType = "favorite"
	Unfavorite    EventType = "unfavorite"
	Follow        EventType = "follow"
)
