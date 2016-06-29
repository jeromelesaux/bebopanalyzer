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

type Point struct {
	Coordinates  string `xml:"coordinates"`
}


type Description struct {
	Data string `xml:",cdata"`
}

//Placemark template
type Placemark struct {
	Name        string `xml:"name"`
	Description Description `xml:"description"`
	StyleUrl string `xml:"styleUrl,omitempty"`
	Point       Point       `xml:"Point,omitempty"`
	LineString []LineString `xml:"LineString,omitempty"`
}


// linestring google template
type LineString struct {
	Extrude      int    `xml:"extrude"`
	AltitudeMode string `xml:"altitudeMode"`
	Coordinates  string `xml:"coordinates"`
}

type Icon struct {
	Href string `xml:"href"`
	Scale float64 `xml:"scale"`
}

type IconStyle struct {
	Id string `xml:"id,attr"`
	Icon Icon `xml:"Icon"`
}

type Style struct {
	Id string `xml:"id,attr"`
	IconStyle IconStyle `xml:"IconStyle"`
}

type Document struct {
	Name      string      `xml:"name"`
	Style Style `xml:"Style,omitempty"`
	Placemark []Placemark `xml:"Placemark"`
}

//Kml template
type Kml struct {
	XMLName   xml.Name `xml:"kml"`
	Namespace string   `xml:"xmlns,attr"`
	AtomNamespace string `xml:"xmlns:atom,attr"`
	GxNamespace string `xml:"xmlns:gx,attr"`
	KmlNamespace string `xml:"xmlns:kml,attr"`
	Document  Document `xml:"Document"`
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
	//kml.XalNamespace = "urn:oasis:names:tc:ciq:xsdschema:xAL:2.0"
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
