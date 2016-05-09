package main

import (
	"bebopanalyzer/configuration"
	"bebopanalyzer/fsmanager"
	"bebopanalyzer/model"
	"bebopanalyzer/routes"
	"fmt"
	"net/http"
	"os"
)

var Conf configuration.AppConfiguration

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Usage : app configuration-file.json bebop-fly.json (optionnal for web server)")
		return
	}

	Conf := configuration.LoadConfiguration(os.Args[1])

	if len(os.Args) == 3 {
		bebopJsonFile := os.Args[2]
		fmt.Println("Parsing " + bebopJsonFile)
		pud := model.Load(bebopJsonFile)
		project := fsmanager.Project{BaseDir: Conf.BasepathStorage, Name: pud.SerialNumber, Data: pud, Date: pud.Date}
		project.PerformAnalyse(pud)
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
