package entities

type pageInfo struct {
	EndCursor string `json:"end_cursor"`
	NextPage  bool   `json:"has_next_page"`
}

type nodeInfo struct {
	ImageURL           string `json:"display_url"`
	ThumbnailURL       string `json:"thumbnail_src"`
	IsVideo            bool   `json:"is_video"`
	Date               int    `json:"date"`
	EdgeMediaToComment struct {
		Count int `json:"count"`
	} `json:"edge_media_to_comment"`
	EdgeLikedBy struct {
		Count int `json:"count"`
	} `json:"edge_liked_by"`
}

type ProfileInfo struct {
	DataInfo struct {
		User UserInfo `json:"user"`
	} `json:"data"`
	Status string `json:"status"`
}

type UserInfo struct {
	EdgeFollowedBy struct {
		Count int `json:"count"`
	} `json:"edge_followed_by"`
	EdgeFollow struct {
		Count int `json:"count"`
	} `json:"edge_follow"`
	Id    string `json:"id"`
	Media struct {
		Edges []struct {
			Node nodeInfo `json:"node"`
		} `json:"edges"`
		PageInfo pageInfo `json:"page_info"`
	} `json:"edge_owner_to_timeline_media"`
}

type ProfileData struct {
	EntryData struct {
		ProfilePage []struct {
			Graphql struct {
				User UserInfo `json:"user"`
			} `json:"graphql"`
		} `json:"ProfilePage"`
	} `json:"entry_data"`
}
