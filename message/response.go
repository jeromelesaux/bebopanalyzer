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
	Place        string `json:"place,omitempty"`
	CsvFile      string `json:"csvFile"`
	KmzFile      string `json:"kmzFile"`
	GpxFile      string `json:"gpxFile"`
	OriginalFile string `json:"originalFile"`
	FlyDuration  string `json:"flyDuration"`
}

type GoogleChartColumns struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

type GoogleRowsValues struct {
	V []interface{} `json:"v"`
}

type GoogleChartRows struct {
	C []GoogleRowsValues `json:"c"`
}

type JsonChartDataResponse struct {
	Value [][]interface{} `json:""`
}
