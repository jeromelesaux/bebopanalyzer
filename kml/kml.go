package kml

import (
	"encoding/xml"
	"strconv"
)

type AltitudeModeKey int

const (
	ClampToGround AltitudeModeKey = iota
	RelativeToGround
	Absolute
	ClampToSeaFloor
	RelativeToSeaFloor
)

var AltitudeMode = [...]string{
	"clampToGround",
	"relativeToGround",
	"absolute",
	"clampToSeaFloor",
	"relativeToSeaFloor",
}

//Placemark template
type Placemark struct {
	Name        string `xml:"kml:name"`
	Description string `xml:"kml:description"`
	//Point       string       `xml:"kml:coordinates"`
	LineString []LineString `xml:"kml:LineString"`
}

// linestring google template
type LineString struct {
	Extrude      int    `xml:"kml:extrude"`
	AltitudeMode string `xml:"kml:altitudeMode"`
	Coordinates  string `xml:"kml:coordinates"`
}

type Document struct {
	Name      string      `xml:"kml:name"`
	Placemark []Placemark `xml:"kml:Placemark"`
}

//Kml template
type Kml struct {
	XMLName       xml.Name `xml:"kml:kml"`
	Namespace     string   `xml:"xmlns,attr"`
	GxNamespace   string   `xml:"xmlns:gx,attr"`
	KmlNamespace  string   `xml:"xmlns:kml,attr"`
	AtomNamespace string   `xml:"xmlns:atom,attr"`
	XalNamespace  string   `xml:"xmlns:xal,attr"`
	Document      Document `xml:"kml:Document"`
}

func NewKML(namespace string, numPlacemarks int) *Kml {
	//Initiate new kml layout
	kml := new(Kml)
	if namespace == "" {
		namespace = "http://www.opengis.net/kml/2.2"
	}
	kml.Namespace = namespace
	kml.AtomNamespace = "http://www.w3.org/2005/Atom"
	kml.GxNamespace = "http://www.google.com/kml/ext/2.2"
	kml.KmlNamespace = "http://www.opengis.net/kml/2.2"
	kml.XalNamespace = "urn:oasis:names:tc:ciq:xsdschema:xAL:2.0"
	kml.Document.Placemark = make([]Placemark, numPlacemarks)
	return kml
}

func (k *Kml) AddPlacemark(placemark Placemark) {
	k.Document.Placemark = append(k.Document.Placemark, placemark)
}

//func (k *Kml) AddPlacemark(name string, desc string, point string) {
//	placemark := Placemark{}
//	placemark.Name = name
//	placemark.Description = desc
//	placemark.Point = point
//	k.Document.Placemark = append(k.Document.Placemark, placemark)
//}

func (p *Placemark) AddLineString(lineString LineString) {
	p.LineString = append(p.LineString, lineString)
}

func (l *LineString) AddCoordinate(longitude float64, latitude float64, altitude float64) {
	//b := bytes.NewBufferString(l.Coordinates)
	//b.WriteString(strconv.FormatFloat(longitude, 'f', 6, 64) + "," + strconv.FormatFloat(latitude, 'f', 6, 64) + "," + strconv.FormatFloat(altitude, 'f', 6, 64) + " ")
	l.Coordinates += strconv.FormatFloat(longitude, 'f', 6, 64) + "," + strconv.FormatFloat(latitude, 'f', 6, 64) + "," + strconv.FormatFloat(altitude, 'f', 6, 64) + " "
}

func (k *Kml) Marshal() ([]byte, error) {
	return xml.MarshalIndent(k, "", "    ")
}
