package domain

import (
	"time"
)

type Article struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ArticleID   string    `json:"articleID" bson:"articleID"`
	TeamID      string    `json:"teamId" bson:"teamId"`
	ClubURL     string    `json:"ClubURL" bson:"clubURL,omitempty"`
	OptaMatchID string    `json:"optaMatchId" bson:"optaMatchId,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Type        []string  `json:"type" bson:"type,omitempty"`
	Teaser      string    `json:"teaser" bson:"teaser,omitempty"`
	Content     string    `json:"content" bson:"content,omitempty"`
	URL         string    `json:"url" bson:"url,omitempty"`
	ImageURL    string    `json:"imageUrl" bson:"imageUrl,omitempty"`
	GalleryURLs []string  `json:"galleryUrls" bson:"galleryUrls,omitempty"`
	VideoURL    string    `json:"videoUrl" bson:"videoUrl,omitempty"`
	BodyText    string    `json:"bodyText" bson:"bodyText,omitempty"`
	Subtitle    string    `json:"subtitle" bson:"subtitle,omitempty"`
	IsPublished bool      `json:"isPublished" bson:"isPublished,omitempty"`
	Published   time.Time `json:"published" bson:"published"`
}

type ArticleRest struct {
	Status string   `json:"status"`
	Data   *Article `json:"data"`
}

func (a *Article) ToRest() *ArticleRest {
	return &ArticleRest{
		Status: "success",
		Data:   a,
	}
}

type Articles struct {
	Total    int64      `json:"total"`
	Articles []*Article `json:"articles"`
}

type ArticlesRest struct {
	Status   string     `json:"status"`
	Data     []*Article `json:"data"`
	Metadata Metadata   `json:"metadata"`
}

type Metadata struct {
	Total int64 `json:"total"`
}

func (a *Articles) ToRest() *ArticlesRest {
	return &ArticlesRest{
		Status:   "success",
		Data:     a.Articles,
		Metadata: Metadata{Total: a.Total},
	}
}
