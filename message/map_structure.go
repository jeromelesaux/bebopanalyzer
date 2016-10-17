package message

type Point struct {
	//Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
}
