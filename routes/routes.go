package routes

import (
	"BebopAnalyzer/configuration"
	"BebopAnalyzer/fsmanager"
	"BebopAnalyzer/model"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

var Conf *configuration.AppConfiguration

func AnalyseFly(w http.ResponseWriter, r *http.Request) {
	log.Print("Calling AnalyseFly.")

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
			log.Println("name is : " + name)
			fileInfo, pud := getReader(part)
			log.Println("Fileinfo : " + fileInfo.Name + " with error :" + fileInfo.Error)
			if fileInfo.Error != "" {
				FileInfoAsResponse(w, fileInfo)
				return
			}

			project := fsmanager.Project{BaseDir: Conf.BasepathStorage, Name: pud.SerialNumber, Data: pud, Date: pud.Date}
			project.CreateBaseFS()
			project.CopyOriginalStruct(pud)
			project.CreateCsvFile(project.GeneratedData + string(filepath.Separator) + fsmanager.CSV_FILE_NAME)
			project.CreateKmlFile(project.GeneratedData + string(filepath.Separator) + fsmanager.GOOGLEEARTH_FILENAME)

			FileInfoAsResponse(w, fileInfo)
			return
		}
	}
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
	for _, zf := range z.File {
		reader, _ := zf.Open()
		json.NewDecoder(reader).Decode(pud)
		return
	}

	return
}

func GetListTree(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDirectories called.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage
	t := p.GetTree()
	JsonAsResponse(w, t)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetFlies called.")
	p := fsmanager.Project{}
	p.BaseDir = Conf.BasepathStorage

	serialNumber := r.URL.Query().Get("serialNumber")
	log.Println("SerialNumber:" + serialNumber)
	flyDate := r.URL.Query().Get("flyDate")
	log.Println("flyDate:" + flyDate)
	isKmz := r.URL.Query().Get("kmz")
	isCsv := r.URL.Query().Get("csv")
	isOrginal := r.URL.Query().Get("original")

	if serialNumber != "" && flyDate != "" {
		if isKmz == "true" {
			r := p.GetKmzFile(serialNumber, flyDate)

			FileAsResponse(w, r, fsmanager.GOOGLEEARTH_FILENAME)
		} else {
			if isCsv == "true" {
				r := p.GetCsvFile(serialNumber, flyDate)

				FileAsResponse(w, r, fsmanager.CSV_FILE_NAME)
			} else {
				if isOrginal == "true" {
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
