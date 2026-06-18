// modexusBot/internal/storage/models.go
package storage

type YoutubeChannel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Handle      string `json:"handle"`
	RSSURL      string `json:"rss_url"`
	AvatarURL   string `json:"avatar_url"`
	AvatarPath  string `json:"avatar_path"`
	LastVideoID string `json:"last_video_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
