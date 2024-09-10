package sys

// 系统运营数据
type SysData struct {
	ApiCount           int64           `json:"ApiCount"`
	AllApiCount        int64           `json:"AllApiCount"`
	WeekApiCount       []DateCount     `json:"WeekApiCount"`
	WeekClientApiCount []ClientIPCount `json:"WeekClientApiCount"`
	ReqApiCount        int64           `json:"ReqApiCount"`
	AllReqApiCount     int64           `json:"AllReqApiCount"`
	SystemCount        int64           `json:"SystemCount"`
	RouterCount        int64           `json:"RouterCount"`
}
type DateCount struct {
	Date  string `json:"Date"`
	Count int64  `json:"Count"`
}
type ClientIPCount struct {
	ClientIP string `json:"ClientIP"`
	Count    int64  `json:"Count"`
}
