package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

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

type Component struct {
	Type       int `json:"type"`
	Components []struct {
		Type     int    `json:"type"`
		Style    int    `json:"style,omitempty"`
		Url      string `json:"url,omitempty"`
		Label    string `json:"label,omitempty"`
		CustomId string `json:"custom_id,omitempty"`
		Options  []struct {
			Id          int64  `json:"id"`
			Label       string `json:"label"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"options,omitempty"`
		Placeholder string `json:"placeholder,omitempty"`
	} `json:"components"`
}

type Embeds struct {
	Content    string      `json:"content"`
	Components []Component `json:"components"`
	Embeds     []Embed     `json:"embeds"`
}

type TicketReport struct {
	ClosedAt    time.Time
	OpenedAt    time.Time
	ChannelID   string
	AuthorID    string
	ClosedBy    string
	ChannelName string
	Messages    []*discordgo.Message
}
