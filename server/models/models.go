package models

type UsdBrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ExchangeApiResponse struct {
	UsdBrl `json:"USDBRL"`
}

type HistoryUsdBrl struct {
	Id        int    `json:"-"`
	Bid       string `json:"bid"`
	Timestamp int64  `json:"timestamp"`
}