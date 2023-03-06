package domain

import (
	"time"
)

type HullArticles struct {
	ClubName            string `xml:"ClubName"`
	ClubWebsiteURL      string `xml:"ClubWebsiteURL"`
	NewsletterNewsItems struct {
		NewsletterNewsItem []HullArticle `xml:"NewsletterNewsItem"`
	} `xml:"NewsletterNewsItems"`
}

type HullArticle struct {
	ArticleURL        string   `xml:"ArticleURL"`
	NewsArticleID     string   `xml:"NewsArticleID"`
	PublishDate       string   `xml:"PublishDate"`
	Taxonomies        []string `xml:"Taxonomies"`
	TeaserText        string   `xml:"TeaserText"`
	Subtitle          string   `xml:"Subtitle"`
	ThumbnailImageURL string   `xml:"ThumbnailImageURL"`
	Title             string   `xml:"Title"`
	BodyText          string   `xml:"BodyText"`
	GalleryImageURLs  []string `xml:"GalleryImageURLs"`
	VideoURL          string   `xml:"VideoURL"`
	OptaMatchID       string   `xml:"OptaMatchId"`
	LastUpdateDate    string   `xml:"LastUpdateDate"`
	IsPublished       bool     `xml:"IsPublished"`
}

type HullArticleInformation struct {
	NewsArticle HullArticle `xml:"NewsArticle"`
}

// ToDomain returns new Article from HullArticle.
func (h *HullArticle) ToDomain(clubName, clubURL, body, subtitle string) *Article {
	publishedDate, _ := time.Parse("2006-01-02 15:04:05", h.PublishDate)

	return &Article{
		ArticleID: h.NewsArticleID,
		// teamID == clubName ?.
		TeamID:      clubName,
		ClubURL:     clubURL,
		OptaMatchID: h.OptaMatchID,
		Title:       h.Title,
		Type:        h.Taxonomies,
		Teaser:      h.TeaserText,
		Content:     body,
		URL:         h.ArticleURL,
		ImageURL:    h.ThumbnailImageURL,
		GalleryURLs: h.GalleryImageURLs,
		VideoURL:    h.VideoURL,
		Subtitle:    subtitle,
		IsPublished: h.IsPublished,
		Published:   publishedDate,
	}
}
