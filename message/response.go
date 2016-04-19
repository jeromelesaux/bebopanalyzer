package message

//type JsonTable struct {
//	FliesDate    []JsonFlyRow `json:"flyDate,omitempty"`
//	SerialNumber []string     `json:"serialNumber,omitempty"`
//}

//type JsonFlyRow struct {
//	FlyDate string `json:"flyDate,omitempty"`
//	CsvFile string `json:"csvFile"`
//	KmzFile string `json:"kmzFile"`
//}

//type JsonFiles struct {
//	CsvFile string `json:"csvFile"`
//	KmzFile string `json:"kmzFile"`
//}

type JsonSerialNumberRow struct {
	SerialNumber string `json:"serialNumber,omitempty"`
}

type JsonDataListResponse struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	FlyDate      string `json:"flyDate,omitempty"`
	CsvFile      string `json:"csvFile"`
	KmzFile      string `json:"kmzFile"`
	OriginalFile string `json:"originalFile"`
}
