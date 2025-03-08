package response

type StatisticsDataRespnse struct {
	Code            int `json:"code"`
	StatisticsItems []*StatisticsItem
}

type StatisticsItem struct {
	Label      string  `json:"label"`
	Value      int     `json:"value"`
	Icon       string  `json:"icon"`
	Color      string  `json:"color"`
	Trend      string  `json:"trend"`
	Percentage float64 `json:"percentage"`
}
