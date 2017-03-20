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

// contains all informations about the files (types csv, gpx, kmz, json)
// and identifications informations
// swagger:response  jsonDataListResponse
type JsonDataListResponse struct {
	// serialnumber of the drone
	SerialNumber string `json:"serialNumber,omitempty"`
	// fly date
	FlyDate string `json:"flyDate,omitempty"`
	// localisation of the fly
	Place string `json:"place,omitempty"`
	// path to get the csv file containing the fly data
	CsvFile string `json:"csvFile"`
	// path to get the kmz file of the fly
	KmzFile string `json:"kmzFile"`
	// path to get the gpx file of the fly
	GpxFile string `json:"gpxFile"`
	// path to get the original json file of the fly
	OriginalFile string `json:"originalFile"`
	// fly duration
	FlyDuration string `json:"flyDuration"`
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
