package response

type StatisticsDataRespnse struct {
	Code            int               `json:"code"`
	StatisticsItems []*StatisticsItem `json:"statisticsItems"`
}

type StatisticsItem struct {
	Label      string  `json:"label"`
	Value      int     `json:"value"`
	Icon       string  `json:"icon"`
	Color      string  `json:"color"`
	Trend      string  `json:"trend"`
	Percentage float64 `json:"percentage"`
}

type UserActionLogsResponse struct {
	Code           int              `json:"code"`
	UserActionLogs []*UserActionLog `json:"userActionLogs"`
}

type UserActionLog struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	UserType string `json:"user_type"`
	Action   string `json:"action"`
	Time     string `json:"time"`
	Module   string `json:"module"`
}

type NewAnimeResponse struct {
	Code   int          `json:"code"`
	Animes []*AnimeItem `json:"animes"`
}

type AnimeItem struct {
	Title      string `json:"title"`      // 动画标题
	Image      string `json:"image"`      // 动画封面图片URL
	UpdateTime string `json:"updateTime"` // 更新时间
}
