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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var CSV_FILE_NAME = "data.csv"
var JSON_FILENAME = "original.json"
var GOOGLEEARTH_FILENAME = "fly.kmz"
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
	err2 := os.MkdirAll(p.RawData, os.ModePerm)
	if err2 != nil {
		fmt.Println("Error with error : " + err2.Error())
	}
	fmt.Println("Creating path : " + p.GeneratedData)
	err3 := os.MkdirAll(p.GeneratedData, os.ModePerm)
	if err3 != nil {
		fmt.Println("Error with error : " + err3.Error())
	}

}

func (p *Project) CopyOriginalStruct(pud *model.PUD) (n int64) {
	output, err := os.Create(p.RawData + string(filepath.Separator) + JSON_FILENAME)
	if err != nil {
		println("Error while creating file : " + err.Error())
	}
	defer output.Close()

	errCopy := json.NewEncoder(output).Encode(pud)
	if errCopy != nil {
		println("Error while copyinh file : " + errCopy.Error())
	}

	return
}

func (p *Project) LoadPUD(serialNumber string, flyDate string) {
	path := p.BaseDir + string(filepath.Separator) + serialNumber + string(filepath.Separator) + flyDate + string(filepath.Separator) + "Raw" + string(filepath.Separator) + JSON_FILENAME

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
	n, errCopy := io.Copy(output, input)
	if errCopy != nil {
		println("Error while copyinh file : " + errCopy.Error())
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
	errSave := w.WriteAll(records)
	if errSave != nil {
		println("Error while copying csv file : " + errSave.Error())
	}
	w.Flush()
}

func (p *Project) CreateKmlFile(filepath string) {
	name := "Fly of the " + p.Data.Date + " for the drone " + p.Data.ProductName + " " + p.Data.SerialNumber
	description := "Fly " + p.Data.Date + " for the drone " + p.Data.ProductName + "<br>Serial number " + p.Data.SerialNumber + "<br>Version " + p.Data.Version + "<br>Hardware version " + p.Data.HardwareVersion + "<br>Software version " + p.Data.SoftwareVersion + "<br>UUID " + p.Data.Uuid + "<br>Number of crashs " + strconv.Itoa(p.Data.Crash) + "<br>Controller application " + p.Data.ControllerApplication + "<br>Controller model " + p.Data.ControllerModel + "<br>Fly duration " + strconv.Itoa(p.Data.TotalRunTime/60000) + " minutes"

	kmlObject := kml.NewKML("", 1)
	placemark := kmlObject.Document.Placemark[0]
	placemark.Name = name
	placemark.Description = description

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
			t.OriginalFile = "./get?serialNumber=" + t.SerialNumber + "&flyDate=" + strings.Replace(t.FlyDate, "+", "%2B", 1) + "&original=true"
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
