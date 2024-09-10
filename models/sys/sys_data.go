package sys

// 系统运营数据
type SysData struct {
	ApiCount           int64           `json:"api_count"`
	AllApiCount        int64           `json:"all_api_count"`
	WeekApiCount       []DateCount     `json:"week_api_count"`
	WeekClientApiCount []ClientIPCount `json:"week_client_api_count"`
	ReqApiCount        int64           `json:"req_api_count"`
	AllReqApiCount     int64           `json:"all_req_api_count"`
	SystemCount        int64           `json:"system_count"`
	RouterCount        int64           `json:"router_count"`
}
type DateCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
type ClientIPCount struct {
	ClientIP string `json:"client_ip"`
	Count    int64  `json:"count"`
}
