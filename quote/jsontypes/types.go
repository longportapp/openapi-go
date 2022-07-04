package jsontypes

type WatchedSecurity struct {
	Symbol    string `json:"symbol"`
	Market    string `json:"market"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	WatchedAt int64  `json:"watched_at,string"`
}

type WatchedGroup struct {
	Id        string             `json:"id"`
	Name      string             `json:"name"`
	Securites []*WatchedSecurity `json:"securities"`
}

type WatchedGroupList struct {
	Groups []*WatchedGroup `json:"groups"`
}
