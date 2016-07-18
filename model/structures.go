package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type PUD struct {
	Version               string          `json:"version"`
	SoftwareVersion       string          `json:"software_version"`
	HardwareVersion       string          `json:"hardware_version"`
	Date                  string          `json:"date"`
	ProductId             int64           `json:"product_id"`
	SerialNumber          string          `json:"serial_number"`
	ProductName           string          `json:"product_name"`
	Uuid                  string          `json:"uuid"`
	RunOrigin             int             `json:"run_origin"`
	ControllerModel       string          `json:"controller_model"`
	ControllerApplication string          `json:"controller_application"`
	ProductStyle          interface{}     `json:"product_style"`
	ProductAccessory      interface{}     `json:"product_accessory"`
	GpsAvailable          bool            `json:"gps_available"`
	GpsLatitude           float64         `json:"gps_latitude"`
	GpsLongitude          float64         `json:"gps_longitude"`
	Crash                 int             `json:"crash"`
	Jump                  int             `json:"jump"`
	RunTime               int             `json:"run_time"`
	TotalRunTime          int             `json:"total_run_time"`
	DetailsHeaders        []string        `json:"details_headers"`
	DetailsData           [][]interface{} `json:"details_data"`
}

type DetailValue struct {
	time                       float64
	battery_level              float64
	controller_gps_longitude   float64
	controller_gps_latitude    float64
	flying_state               float64
	alert_state                float64
	wifi_signal                float64
	product_gps_available      bool
	product_gps_longitude      float64
	product_gps_latitude       float64
	product_gps_position_error float64
	product_gps_sv_number      int64
	speed_vx                   float64
	speed_vy                   float64
	speed_vz                   float64
	angle_phi                  float64
	angle_theta                float64
	angle_psi                  float64
	altitude                   float64
	flip_type                  float64
	speed                      float64
}

func (pud *PUD) IndexForKey(keySearch string) int {
	value := 0
	for i := 0; i < len(pud.DetailsHeaders); i++ {
		if pud.DetailsHeaders[i] == keySearch {
			value = i
			break
		}
	}
	return value
}

func (pud *PUD) TimeAt(index int) float64 {
	value := pud.IndexForKey("time")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) BatteryLevelAt(index int) float64 {
	value := pud.IndexForKey("battery_level")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) ControllerGpsLongitudelAt(index int) float64 {
	value := pud.IndexForKey("controller_gps_longitude")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) FlyingStateAt(index int) float64 {
	value := pud.IndexForKey("flying_state")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) ProductGpsAvailableAt(index int) bool {
	value := pud.IndexForKey("product_gps_available")
	return pud.DetailsData[index][value].(bool)
}

func (pud *PUD) AlertStateAt(index int) float64 {
	value := pud.IndexForKey("alert_state")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) WifiSignalAt(index int) float64 {
	value := pud.IndexForKey("wifi_signal")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) ProductGpsLongitudeAt(index int) float64 {
	value := pud.IndexForKey("product_gps_longitude")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) ProductGpsLatidudeAt(index int) float64 {
	value := pud.IndexForKey("product_gps_latitude")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) ProductGpsPositionErrorAt(index int) float64 {
	value := pud.IndexForKey("product_gps_position_error")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) SpeedVxAt(index int) float64 {
	value := pud.IndexForKey("speed_vx")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) SpeedVyAt(index int) float64 {
	value := pud.IndexForKey("speed_vy")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) SpeedVzAt(index int) float64 {
	value := pud.IndexForKey("speed_vz")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) AnglePhiAt(index int) float64 {
	value := pud.IndexForKey("angle_phi")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) AngleThetaAt(index int) float64 {
	value := pud.IndexForKey("angle_theta")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) AnglePsiAt(index int) float64 {
	value := pud.IndexForKey("angle_psi")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) AltitudeAt(index int) float64 {
	value := pud.IndexForKey("altitude")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) FlipTypeAt(index int) float64 {
	value := pud.IndexForKey("flip_type")
	return pud.DetailsData[index][value].(float64)
}

func (pud *PUD) SpeedAt(index int) float64 {
	value := pud.IndexForKey("speed")
	return pud.DetailsData[index][value].(float64)
}

func Load(input string) *PUD {
	pud := &PUD{}
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error while opening file ", input, err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(pud)
	if err != nil {
		fmt.Println("error:", err)
	}
	return pud
}

func (p *PUD) Csv() [][]string {
	length := len(p.DetailsData)
	var records [][]string
	records = append(records, []string{"time", "battery_level", "controller_gps_longitude", "flying_state", "alert_state", "wifi_signal", "product_gps_available", "product_gps_longitude", "product_gps_latitude", "product_gps_position_error", "speed_vx", "speed_vy", "speed_vz", "angle_phi", "angle_theta", "angle_psi", "altitude", "flip_type", "speed"})
	for i := 0; i < length; i++ {
		time := strconv.FormatFloat(p.TimeAt(i), 'f', 6, 64)
		batteryLevel := strconv.FormatFloat(p.BatteryLevelAt(i), 'f', 6, 64)
		controllerGpsLong := strconv.FormatFloat(p.ControllerGpsLongitudelAt(i), 'f', 6, 64)
		flyingState := strconv.FormatFloat(p.FlyingStateAt(i), 'f', 6, 64)
		alertState := strconv.FormatFloat(p.AlertStateAt(i), 'f', 6, 64)
		wifiSignal := strconv.FormatFloat(p.WifiSignalAt(i), 'f', 6, 64)
		productGpsAvailable := strconv.FormatBool(p.ProductGpsAvailableAt(i))
		productGpsLongitude := strconv.FormatFloat(p.ProductGpsLongitudeAt(i), 'f', 6, 64)
		productGpsLatitude := strconv.FormatFloat(p.ProductGpsLatidudeAt(i), 'f', 6, 64)
		productGpsPositionError := strconv.FormatFloat(p.ProductGpsPositionErrorAt(i), 'f', 6, 64)
		speedVx := strconv.FormatFloat(p.SpeedVxAt(i), 'f', 6, 64)
		speedVy := strconv.FormatFloat(p.SpeedVyAt(i), 'f', 6, 64)
		speedVz := strconv.FormatFloat(p.SpeedVzAt(i), 'f', 6, 64)
		anglePhi := strconv.FormatFloat(p.AnglePhiAt(i), 'f', 6, 64)
		angleTheta := strconv.FormatFloat(p.AngleThetaAt(i), 'f', 6, 64)
		anglePsi := strconv.FormatFloat(p.AnglePsiAt(i), 'f', 6, 64)
		altitude := strconv.FormatFloat(p.AltitudeAt(i), 'f', 6, 64)
		flipType := strconv.FormatFloat(p.FlipTypeAt(i), 'f', 6, 64)
		speed := strconv.FormatFloat(p.SpeedAt(i), 'f', 6, 64)
		records = append(records, []string{time, batteryLevel, controllerGpsLong, flyingState, alertState, wifiSignal, productGpsAvailable, productGpsLongitude, productGpsLatitude, productGpsPositionError, speedVx, speedVy, speedVz, anglePhi, angleTheta, anglePsi, altitude, flipType, speed})
	}
	return records
}
