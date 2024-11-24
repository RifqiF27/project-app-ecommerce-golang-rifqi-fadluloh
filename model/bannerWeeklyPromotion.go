package model

type BannerWeeklyPromotionRecomment struct {
	ID       int    `json:"id"`
	Image    string `json:"image"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	PathPage string `json:"path_page,omitempty"`
}
