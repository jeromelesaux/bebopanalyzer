package main

import (
	"BebopAnalyzer/configuration"
	"BebopAnalyzer/fsmanager"
	"BebopAnalyzer/model"
	"BebopAnalyzer/routes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

var Conf configuration.AppConfiguration

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage : app configuration-file.json bebop-fly.json")
		return
	}

	Conf := configuration.LoadConfiguration(os.Args[1])

	if len(os.Args) == 3 {
		bebopJsonFile := os.Args[2]
		fmt.Println("Parsing " + bebopJsonFile)
		pud := &model.PUD{}
		pud = pud.Load(bebopJsonFile)
		project := fsmanager.Project{BaseDir: Conf.BasepathStorage, Name: pud.SerialNumber, Data: pud, Date: pud.Date}
		project.CreateBaseFS()
		project.CopyOriginalFile(bebopJsonFile)
		project.CreateCsvFile(project.GeneratedData + string(filepath.Separator) + fsmanager.CSV_FILE_NAME)
		project.CreateKmlFile(project.GeneratedData + string(filepath.Separator) + fsmanager.GOOGLEEARTH_FILENAME)
		fmt.Println("job done and generated new file at " + project.GeneratedData)
	} else {
		fmt.Println("Starting server web at port " + Conf.HttpPort)
		// gestion des routes http
		routes.Conf = Conf
		http.HandleFunc("/import/", routes.AnalyseFly)
		http.HandleFunc("/get", routes.GetFile)
		http.HandleFunc("/list", routes.GetListTree)
		http.HandleFunc("/chart", routes.GetChart)
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./resources"))))
		http.ListenAndServe(Conf.HttpPort, nil)
	}

}
