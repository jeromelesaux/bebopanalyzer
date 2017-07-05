package kml

import (
	"encoding/xml"
	"fmt"
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
	AltitudeMode string `xml:"altitudeMode,omitempty"`
	Coordinates  string `xml:"coordinates"`
}

func NewCoordinates(longitude, latitude, altitude float64) string {
	return fmt.Sprintf("%f,%f,%f", longitude, latitude, altitude)
}
func NewCoordinatesOffset(longitude, latitude float64) string {
	return fmt.Sprintf("%f,%f,(AltitudeAboveSeaLevelOffset+1)", longitude, latitude)
}

type Description struct {
	Data string `xml:",cdata"`
}

//Placemark template
type Placemark struct {
	Name        string       `xml:"name"`
	Description Description  `xml:"description"`
	StyleUrl    string       `xml:"styleUrl"`
	Style       Style        `xml:"Style,omitempty"`
	Point       Point        `xml:"Point,omitempty"`
	LineString  []LineString `xml:"LineString,omitempty"`
}

// linestring google template
type LineString struct {
	Extrude      int    `xml:"extrude"`
	AltitudeMode string `xml:"altitudeMode"`
	Coordinates  string `xml:"coordinates"`
	Tessellate   int    `xml:"tessellate"`
}

type Icon struct {
	Href string `xml:"href,omitempty"`
}

type IconStyle struct {
	Id    string  `xml:"id,attr,omitempty"`
	Icon  Icon    `xml:"Icon,omitempty"`
	Color string  `xml:"color,omitempty"`
	Scale float64 `xml:"scale,omitempty"`
}

type LabelStyle struct {
	Color string  `xml:"color,omitempty"`
	Scale float64 `xml:"scale,omitempty"`
}
type LineStyle struct {
	Color string `xml:"color,omitempty"`
	Width int    `xml:"width,omitempty"`
}

type Style struct {
	Id         string     `xml:"id,attr,omitempty"`
	IconStyle  IconStyle  `xml:"IconStyle,omitempty"`
	LabelStyle LabelStyle `xml:"LabelStyle,omitempty"`
	LineStyle  LineStyle  `xml:"LineStyle,omitempty"`
}

type Document struct {
	Name      string      `xml:"name"`
	Style     Style       `xml:"Style"`
	Placemark []Placemark `xml:"Placemark"`
}

//Kml template
type Kml struct {
	XMLName       xml.Name `xml:"kml"`
	Namespace     string   `xml:"xmlns,attr"`
	AtomNamespace string   `xml:"xmlns:atom,attr"`
	GxNamespace   string   `xml:"xmlns:gx,attr"`
	KmlNamespace  string   `xml:"xmlns:kml,attr"`
	Document      Document `xml:"Document"`
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
