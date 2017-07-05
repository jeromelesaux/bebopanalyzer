// Package routes bebopanalyzer API
// the purpose of this application is to parse bebop metadata fly file
// it will provided mutliple format files to display analyze yours flies
//
// Terms Of Service:
// bebop metadata flies json zipped files manager
//
// Title: bebopanalyzer
// Schemes: https
// Host: localhost
// BasePath: /
// Version: 1.9
// Licence: MIT http://opensource.org/licenses/MIT
// Contact: jeromelesaux@gmail.com
//
// Consumes:
// - application/json
// - multipart/form-data
//
// Produces:
// - application/json
//
// swagger:meta
package routes

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jeromelesaux/bebopanalyzer/configuration"
	"github.com/jeromelesaux/bebopanalyzer/fsmanager"
	"github.com/jeromelesaux/bebopanalyzer/model"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var Conf *configuration.AppConfiguration

func importFlyWG(pud *model.PUD, fis *FileInfos, fileInfo *FileInfo, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fis.Files = append(fis.Files, fileInfo)
		if fileInfo.Error == "" {
			importFly(pud)
		}
		wg.Done()
	}()

}

// swagger:parameters import
type FileImported struct {
	// required: true
	File []byte
}

// AnalyseFly swagger:route POST /import  import
//
// Import JSON fly metadata
//
//	Consumes:
// 	- multipart/form-data
//
//	Produces :
//	- application/json
//
//	Schemes: https
//
// 	Security:
//
//	Responses :
//	default: fileInfos
func AnalyseFly(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	log.Print("Calling AnalyseFly.")
	fis := &FileInfos{Files: []*FileInfo{}}

	defer func() {
		JsonAsResponse(w, fis)
	}()

	w.Header().Set("Content-Type", "application/javascript")
	mr, err := r.MultipartReader()
	if err != nil {
		//msg := "500 Internal Server Error: " + err.Error()
		//http.Error(w, msg, http.StatusInternalServerError)
		//fmt.Println(err.Error())
		return
	}
	r.Form, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		//msg := "500 Internal Server Error: " + err.Error()
		//http.Error(w, msg, http.StatusInternalServerError)
		//fmt.Println(err.Error())
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			//msg := "500 Internal Server Error: " + err.Error()
			//http.Error(w, msg, http.StatusInternalServerError)
			//fmt.Println(err.Error())
			return
		}
		if name := part.FormName(); name != "" {
			fileInfo, pud := getReader(part)
			log.Println("Fileinfo : " + fileInfo.Name + " with error :" + fileInfo.Error)
			importFlyWG(pud, fis, fileInfo, &wg)
			continue
		}
	}
	wg.Wait()
}

func importFly(pud *model.PUD) {
	project := fsmanager.Project{BaseDir: Conf.BasepathStorage, Name: pud.SerialNumber, Data: pud, Date: pud.Date}
	project.PerformAnalyse(pud)
	return
}

func getReader(p *multipart.Part) (fi *FileInfo, pud *model.PUD) {
	fi = &FileInfo{
		Name: p.FileName(),
	}
	// Validate file type
	if !strings.HasSuffix(strings.ToLower(fi.Name), ".zip") {
		fi.Error = "Ceci n'est pas un fichier ZIP."
		return
	}
	pud = &model.PUD{}

	b, err := ioutil.ReadAll(p)
	if err != nil {
		log.Println("ERror while copying part " + err.Error())
		return
	}
	br := bytes.NewReader(b)
	z, _ := zip.NewReader(br, int64(len(b)))
	fi.Size = int64(len(b))
	fmt.Println(fi)
	for _, zf := range z.File {
		reader, _ := zf.Open()
		err := json.NewDecoder(reader).Decode(pud)
		if err != nil {
			fmt.Println("Error while unmarshalling file with error " + err.Error())
		}
		return
	}

	return
}

// GetListTree swagger:route GET /list getListTree
//
// Get file for type type
//
//	Consumes:
// 	- application/json
//
//	Produces :
//	- application/json
//
//	Schemes: https
//
// 	Security:
//
//	Responses :
//	default: jsonDataListResponse
func GetListTree(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDirectories called.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage
	t := p.GetTree()
	JsonAsResponse(w, t)
}

// swagger:parameters getFile
type GetFileParameters struct {

	// serial number of the drone
	// required: true
	SerialNumber string

	// fly date of instance
	// required: true
	FlyDate string

	// type accepted csv, kmz, gpx, original
	// required: true
	Type string
}

// swagger:response getFileResponse
type GetFileResponse struct {
	File []byte
}

// GetFile swagger:route GET /get/{serialNumber}/{flyDate}/{type} serialNumber flyDate type getFile
//
// Get file for type type
//
//	Consumes:
// 	- application/json
//
//	Produces :
//	- multipart/form-data
//
//	Schemes: https
//
// 	Security:
//
//	Responses :
//	default: getFileResponse
func GetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFlies called.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage

	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]
	//serialNumber := r.URL.Query().Get("serialNumber")
	log.Println("SerialNumber:" + serialNumber)
	//flyDate := r.URL.Query().Get("flyDate")
	flyDate := vars["flyDate"]
	log.Println("flyDate:" + flyDate)
	//isKmz := r.URL.Query().Get("kmz")
	//isCsv := r.URL.Query().Get("csv")
	//isGpx := r.URL.Query().Get("gpx")
	//isOrginal := r.URL.Query().Get("original")
	typeFile := vars["type"]
	if serialNumber != "" && flyDate != "" {
		if typeFile == "kmz" {
			r := p.GetKmzFile(serialNumber, flyDate)
			FileAsResponse(w, r, fsmanager.GOOGLEEARTH_FILENAME)
		} else {
			if typeFile == "gpx" {
				p.LoadPUD(serialNumber, flyDate)
				r := p.GetGPXData(serialNumber, flyDate)

				FileAsResponse(w, r, fsmanager.GPX_FILENAME)
			} else {
				if typeFile == "csv" {
					r := p.GetCsvFile(serialNumber, flyDate)
					FileAsResponse(w, r, fsmanager.CSV_FILE_NAME)
				} else {
					if typeFile == "original" {
						r := p.GetOriginalFile(serialNumber, flyDate)

						FileAsResponse(w, r, fsmanager.JSON_FILENAME)
					} else {
						msg := "500 Internal Server Error: "
						http.Error(w, msg, http.StatusInternalServerError)
						fmt.Println("Cannot get file content")
						return
					}
				}
			}
		}

	}

}

// swagger:parameters getChart
type GetChartParameters struct {

	// serial number of the drone
	// required: true
	SerialNumber string

	// fly date of instance
	// required: true
	FlyDate string
}

// GetChart swagger:route GET /chart/{serialNumber}/{flyDate} serialNumber flyDate getChart
//
// Get file for type type
//
//	Consumes:
// 	- application/json
//
//	Produces :
//	- application/json
//
//	Schemes: https
//
// 	Security:
//
//	Responses :
//	default:
func GetChart(w http.ResponseWriter, r *http.Request) {
	log.Println("GetChart calling.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage
	vars := mux.Vars(r)
	serialNumber := vars["serialNumber"]
	//serialNumber := r.URL.Query().Get("serialNumber")
	log.Println("SerialNumber:" + serialNumber)
	//flyDate := r.URL.Query().Get("flyDate")
	flyDate := vars["flyDate"]
	p.LoadPUD(serialNumber, flyDate)
	data := p.GetChartData(serialNumber, flyDate)
	JsonAsResponse(w, data)
}

// swagger:parameters getMaps
type GetMapsParameters struct {

	// serial number of the drone
	// required: true
	SerialNumber string

	// fly date of instance
	// required: true
	FlyDate string
}

// GetMaps swagger:route GET /displayFly/{serialNumber}/{flyDate} serialNumber flyDate getMaps
//
// Get file for type type
//
//	Consumes:
// 	- application/json
//
//	Produces :
//	- application/json
//
//	Schemes: http
//
// 	Security:
//
//	Responses :
//	default: point
func GetMaps(w http.ResponseWriter, r *http.Request) {
	log.Println("GetMaps calling.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage
	vars := mux.Vars(r)
	// swagger:parameters serialNumber
	serialNumber := vars["serialNumber"]
	//serialNumber := r.URL.Query().Get("serialNumber")
	log.Println("SerialNumber:" + serialNumber)
	//flyDate := r.URL.Query().Get("flyDate")
	// swagger:parameters flydate
	flyDate := vars["flyDate"]
	log.Println("flyDate:" + flyDate)
	p.LoadPUD(serialNumber, flyDate)
	data := p.GetMapsData(serialNumber, flyDate)
	JsonAsResponse(w, data)
}

func RebuildData(w http.ResponseWriter, r *http.Request) {
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage
	go func() {
		p.RebuildDataFiles(Conf)
	}()
	JsonAsResponse(w, "rebuilding all your data")
}

func FileAsResponse(w http.ResponseWriter, streamBytes []byte, filename string) {
	log.Println("Writing file.")
	b := bytes.NewBuffer(streamBytes)

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+filename+"\"")
	if _, err := b.WriteTo(w); err != nil {
		log.Println(w, "%s", err.Error())
	}

}

func FileInfoAsResponse(w http.ResponseWriter, fileInfo *FileInfo) {
	fis := &FileInfos{Files: []*FileInfo{fileInfo}}
	JsonAsResponse(w, fis)
}

func JsonAsResponse(w http.ResponseWriter, o interface{}) {
	js, err := json.Marshal(o)
	if err != nil {
		log.Println("Error while marshalling  object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application-json")
	w.Write(js)
}
