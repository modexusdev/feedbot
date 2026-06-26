// modexusBot/internal/storage/models.go
package storage

// YoutubeChannel represents a tracked YouTube channel.
type YoutubeChannel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Handle      string `json:"handle"`
	RSSURL      string `json:"rss_url"`
	AvatarURL   string `json:"avatar_url"`
	LastVideoID string `json:"last_video_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// WeatherLocation represents a tracked weather location.
type WeatherLocation struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
