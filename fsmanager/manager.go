package fsmanager

import (
	"archive/zip"
	"bebopanalyzer/kml"
	"bebopanalyzer/message"
	"bebopanalyzer/model"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/metakeule/fmtdate"
	"github.com/ptrv/go-gpx"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var CSV_FILE_NAME = "data.csv"
var JSON_FILENAME = "original.json"
var GOOGLEEARTH_FILENAME = "fly.kmz"
var GPX_FILENAME = "fly.gpx"
var GOOGLEEARTH_INTERNAL_FILENAME = "doc.kml"

//var BASEDIR_NAME = "Data"

type SortByModified []os.FileInfo

func (f SortByModified) Len() int           { return len(f) }
func (f SortByModified) Less(i, j int) bool { return f[i].ModTime().Unix() > f[j].ModTime().Unix() }
func (f SortByModified) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type Project struct {
	BaseDir       string
	Name          string
	Date          string
	RawData       string
	GeneratedData string
	Data          *model.PUD
}

func (p *Project) PerformAnalyse(pud *model.PUD) {
	p.CreateBaseFS()
	p.CopyOriginalStruct(pud)
	p.CreateCsvFile(p.GeneratedData + string(filepath.Separator) + CSV_FILE_NAME)
	p.CreateKmlFile(p.GeneratedData + string(filepath.Separator) + GOOGLEEARTH_FILENAME)
}

func (p *Project) CreateBaseFS() {

	if p.Data != nil {
		p.RawData = "." + string(filepath.Separator) + p.BaseDir + string(filepath.Separator) + p.Name + string(filepath.Separator) + p.Date + string(filepath.Separator) + "Raw"
		p.GeneratedData = "." + string(filepath.Separator) + p.BaseDir + string(filepath.Separator) + p.Name + string(filepath.Separator) + p.Date + string(filepath.Separator) + "Generated"
	} else {
		fmt.Println("Data is nil!!!")
	}
	fmt.Println("Creating path : " + p.BaseDir)
	err := os.MkdirAll(p.BaseDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error with error : " + err.Error())
	}
	fmt.Println("Creating path : " + p.RawData)
	err = os.MkdirAll(p.RawData, os.ModePerm)
	if err != nil {
		fmt.Println("Error with error : " + err.Error())
	}
	fmt.Println("Creating path : " + p.GeneratedData)
	err = os.MkdirAll(p.GeneratedData, os.ModePerm)
	if err != nil {
		fmt.Println("Error with error : " + err.Error())
	}

}

func (p *Project) CopyOriginalStruct(pud *model.PUD) (n int64) {
	output, err := os.Create(p.RawData + string(filepath.Separator) + JSON_FILENAME)
	if err != nil {
		println("Error while creating file : " + err.Error())
	}
	defer output.Close()

	err = json.NewEncoder(output).Encode(pud)
	if err != nil {
		println("Error while copyinh file : " + err.Error())
	}

	return
}

func (p *Project) LoadPUD(serialNumber string, flyDate string) {
	path := p.BaseDir + string(filepath.Separator) + serialNumber + string(filepath.Separator) + flyDate + string(filepath.Separator) + "Raw" + string(filepath.Separator) + JSON_FILENAME
	log.Println("path to search :" + path)
	fReader, err := os.Open(path)
	if err != nil {
		log.Println("Error while reading file " + path + " with error " + err.Error())
		return
	}
	err = json.NewDecoder(fReader).Decode(&p.Data)
	if err != nil {
		log.Println("Error while decode from file " + path + " with error " + err.Error())
		return
	}
	return
}

func (p *Project) CopyOriginalFile(file string) int64 {
	input, errOpen := os.Open(file)
	if errOpen != nil {
		fmt.Println("Error while opening file ", file, errOpen.Error())
	}
	defer input.Close()

	output, err := os.Create(p.RawData + string(filepath.Separator) + JSON_FILENAME)
	if err != nil {
		println("Error while creating file : " + err.Error())
	}
	defer output.Close()
	n, err := io.Copy(output, input)
	if err != nil {
		println("Error while copyinh file : " + err.Error())
	}
	fmt.Println("Wrote " + strconv.FormatInt(n, 10) + " from file : " + file)
	return n
}

func (p *Project) CreateCsvFile(file string) {
	records := p.Data.Csv()
	csvHandler, err := os.Create(file)
	if err != nil {
		println("Error while creating csv file : " + err.Error())
	}
	defer csvHandler.Close()
	w := csv.NewWriter(csvHandler)
	w.Comma = ';'
	err = w.WriteAll(records)
	if err != nil {
		println("Error while copying csv file : " + err.Error())
	}
	w.Flush()
}

func (p *Project) CreateKmlFile(filepath string) {
	name := "Fly of the " + p.Data.Date + " for the drone " + p.Data.ProductName + " " + p.Data.SerialNumber
	description := "Fly " + p.Data.Date + " for the drone " + p.Data.ProductName + "<br>Serial number " + p.Data.SerialNumber + "<br>Version " + p.Data.Version + "<br>Hardware version " + p.Data.HardwareVersion + "<br>Software version " + p.Data.SoftwareVersion + "<br>UUID " + p.Data.Uuid + "<br>Number of crashs " + strconv.Itoa(p.Data.Crash) + "<br>Controller application " + p.Data.ControllerApplication + "<br>Controller model " + p.Data.ControllerModel + "<br>Fly duration " + strconv.Itoa(p.Data.TotalRunTime/60000) + " minutes"

	kmlObject := kml.NewKML("", 1)
	kmlObject.Document.Name = name
	placemark := kmlObject.Document.Placemark[0]
	placemark.Name = name
	placemark.Description.Data = description

	lineString := kml.LineString{}
	lineString.AltitudeMode = kml.AltitudeMode[kml.RelativeToGround]
	lineString.Extrude = 1

	for i := 0; i < len(p.Data.DetailsData); i++ {
		gpsAvailable := p.Data.ProductGpsAvailableAt(i)

		if gpsAvailable {
			longitude := p.Data.ProductGpsLongitudeAt(i)
			latitude := p.Data.ProductGpsLatidudeAt(i)
			altitude := p.Data.AltitudeAt(i) / 1000
			if latitude != 500. {
				lineString.AddCoordinate(longitude, latitude, altitude)
			}
		}
	}
	placemark.AddLineString(lineString)
	kmlObject.AddPlacemark(placemark)

	content, errMarshalling := xml.Marshal(kmlObject)
	if errMarshalling != nil {
		println("Error while marshalling kml : " + errMarshalling.Error())
	}

	b, _ := os.Create(filepath)
	z := zip.NewWriter(b)
	kmlContent := string(content)
	var files = []struct{ Name, Body string }{{GOOGLEEARTH_INTERNAL_FILENAME, kmlContent}}
	for _, file := range files {
		f, err := z.Create(file.Name)
		if err != nil {
			println("Error while creating  kmz file : " + err.Error())
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			println("Error while creating  kmz file : " + err.Error())
		}
	}
	z.Close()
}

func (p *Project) ListAllDrones() []message.JsonSerialNumberRow {
	r := []message.JsonSerialNumberRow{}
	f, err := ioutil.ReadDir(p.BaseDir)

	if err != nil {
		log.Println("Error while scanning directory " + p.BaseDir + " " + err.Error())
	}
	for _, file := range f {
		el := message.JsonSerialNumberRow{file.Name()}
		r = append(r, el)
	}
	return r
}

func (p *Project) GetTree() []message.JsonDataListResponse {
	r := []message.JsonDataListResponse{}
	root := p.BaseDir

	dir, err := ioutil.ReadDir(root)
	if err != nil {
		log.Println("Error while reading basedir :" + root + " " + err.Error())
	}

	for _, f := range SortByModified(dir) {
		subDir, err := ioutil.ReadDir(root + string(filepath.Separator) + f.Name())
		if err != nil {
			log.Println("Error while reading basedir :" + f.Name() + " " + err.Error())
		}
		for _, subF := range subDir {
			t := message.JsonDataListResponse{}
			t.SerialNumber = f.Name()
			t.FlyDate = subF.Name()
			t.CsvFile = "./get?serialNumber=" + t.SerialNumber + "&flyDate=" + strings.Replace(t.FlyDate, "+", "%2B", 1) + "&csv=true"
			t.KmzFile = "./get?serialNumber=" + t.SerialNumber + "&flyDate=" + strings.Replace(t.FlyDate, "+", "%2B", 1) + "&kmz=true"
			t.GpxFile = "./get?serialNumber=" + t.SerialNumber + "&flyDate=" + strings.Replace(t.FlyDate, "+", "%2B", 1) + "&gpx=true"
			t.OriginalFile = "./get?serialNumber=" + t.SerialNumber + "&flyDate=" + strings.Replace(t.FlyDate, "+", "%2B", 1) + "&original=true"
			p := model.Load(root + string(filepath.Separator) + f.Name() + string(filepath.Separator) + t.FlyDate + string(filepath.Separator) + "Raw" + string(filepath.Separator) + JSON_FILENAME)
			t.FlyDuration = fmt.Sprintf("%.2f", (float32(p.TotalRunTime) / 60000))
			r = append(r, t)
		}
	}
	return r
}

func (p *Project) GetKmzFile(serialNumber string, flyDate string) []byte {

	path := p.BaseDir + string(filepath.Separator) + serialNumber + string(filepath.Separator) + flyDate + string(filepath.Separator) + "Generated" + string(filepath.Separator) + GOOGLEEARTH_FILENAME
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Error while reading file :" + path + " " + err.Error())
	}

	return f
}

func (p *Project) GetCsvFile(serialNumber string, flyDate string) []byte {

	path := p.BaseDir + string(filepath.Separator) + serialNumber + string(filepath.Separator) + flyDate + string(filepath.Separator) + "Generated" + string(filepath.Separator) + CSV_FILE_NAME
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Error while reading file :" + path + " " + err.Error())
	}

	return f
}

func (p *Project) GetOriginalFile(serialNumber string, flyDate string) []byte {

	path := p.BaseDir + string(filepath.Separator) + serialNumber + string(filepath.Separator) + flyDate + string(filepath.Separator) + "Raw" + string(filepath.Separator) + JSON_FILENAME
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("Error while reading file :" + path + " " + err.Error())
	}

	return f
}

func (p *Project) GetChartData(serialNumber string, flyDate string) [][]interface{} {
	m := make([][]interface{}, 0)

	for i := 0; i < len(p.Data.DetailsData); i++ {
		gpsAvailable := p.Data.ProductGpsAvailableAt(i)

		if gpsAvailable {
			r := make([]interface{}, 4)
			time := p.Data.TimeAt(i) / 60000
			speed := p.Data.SpeedAt(i)
			altitude := p.Data.AltitudeAt(i) / 1000
			batteryLevel := p.Data.BatteryLevelAt(i)
			r[0] = time
			r[1] = speed
			r[2] = altitude
			r[3] = batteryLevel
			//c := message.JsonChartDataResponse{time, batteryLevel, altitude, speed}
			//log.Println("time:" + strconv.FormatFloat(time, 'f', 6, 10))
			m = append(m, r)
		}
	}
	return m
}

func (p *Project) GetMapsData(serialNumber string, flyDate string) []message.Point {
	m := []message.Point{}
	log.Println("Entered in GetMapsData and data length is " + strconv.Itoa(len(p.Data.DetailsData)))
	for i, _ := range p.Data.DetailsData {
		gpsAvailable := p.Data.ProductGpsAvailableAt(i)
		//log.Println("GpsIsavailable :" + strconv.FormatBool(gpsAvailable))
		if gpsAvailable {
			latitude := p.Data.ProductGpsLatidudeAt(i)
			longitude := p.Data.ProductGpsLongitudeAt(i)
			//time := p.Data.TimeAt(i) / 60000
			//name := strconv.FormatFloat(time, 'f', 8, 64)
			//log.Println("latitude:" + strconv.FormatFloat(latitude, 'f', 10, 64) + "&longitude:" + strconv.FormatFloat(longitude, 'f', 10, 64))
			if latitude != 500. && longitude != 500. {
				point := message.Point{Description: "", Latitude: latitude, Longitude: longitude}
				m = append(m, point)
			}
		}
	}

	return m
}

func (p *Project) GetGPXData(serialNumber string, flyDate string) []byte {
	name := "Fly " + p.Data.Date + " for the drone " + p.Data.ProductName + " Serial number " + p.Data.SerialNumber + " Version " + p.Data.Version + " Hardware version " + p.Data.HardwareVersion + " Software version " + p.Data.SoftwareVersion + " UUID " + p.Data.Uuid + " Number of crashs " + strconv.Itoa(p.Data.Crash) + " Controller application " + p.Data.ControllerApplication + " Controller model " + p.Data.ControllerModel + " Fly duration " + strconv.Itoa(p.Data.TotalRunTime/60000) + " minutes"
	gpxObject := gpx.NewGpx()
	person := &gpx.Person{Name: name}
	gpxObject.Metadata = &gpx.Metadata{}
	gpxObject.Metadata.Author = person

	trk := gpx.Trk{}
	trk.Name = p.Data.Date
	trkSeg := gpx.Trkseg{}

	startTime, err := fmtdate.Parse("YYYY-MM-DDThhmmss+ssss", p.Data.Date)
	if err != nil {
		fmt.Println(err.Error())
	}

	for key, _ := range p.Data.DetailsData {
		if p.Data.ProductGpsLatidudeAt(key) != 500. && p.Data.ProductGpsLongitudeAt(key) != 500. {
			secondes := p.Data.TimeAt(key) / 60000
			when := time.Duration(secondes) * time.Second
			startTime = startTime.Add(when)

			trkpt := gpx.Wpt{Lat: p.Data.ProductGpsLatidudeAt(key),
				Lon:       p.Data.ProductGpsLongitudeAt(key),
				Ele:       p.Data.AltitudeAt(key) / 1000,
				Timestamp: startTime.Format(time.RFC3339)}

			trkSeg.Waypoints = append(trkSeg.Waypoints, trkpt)
		}

	}
	trk.Segments = append(trk.Segments, trkSeg)
	gpxObject.Tracks = append(gpxObject.Tracks, trk)

	bytes, err := xml.Marshal(gpxObject)
	if err != nil {
		fmt.Println(err.Error())
		return bytes
	}
	return bytes
}
