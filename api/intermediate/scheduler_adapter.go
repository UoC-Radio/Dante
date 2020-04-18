package ws

import (
	"bytes"
	"encoding/xml"
	"time"
)

type Schedule struct {
	XMLName xml.Name `xml:"WeekSchedule"`
	Mon     Day      `xml:" Mon"`
	Tue     Day      `xml:" Tue"`
	Wed     Day      `xml:" Wed"`
	Thu     Day      `xml:" Thu"`
	Fri     Day      `xml:" Fri"`
	Sat     Day      `xml:" Sat"`
	Sun     Day      `xml:" Sun"`
}

type Day struct {
	Zone []Zone `xml:" Zone"`
}

type Fader struct {
	FadeInDurationSecs  int     `xml:" FadeInDurationSecs,omitempty"`
	FadeOutDurationSecs int     `xml:" FadeOutDurationSecs,omitempty"`
	MinLevel            float32 `xml:" MinLevel,omitempty"`
	MaxLevel            float32 `xml:" MaxLevel,omitempty"`
}

type IntermediatePlaylist struct {
	Path              string `xml:" Path"`
	Shuffle           bool   `xml:" Shuffle"`
	Fader             Fader  `xml:" Fader,omitempty"`
	SchedIntervalMins int    `xml:" SchedIntervalMins"`
	NumSchedItems     int    `xml:" NumSchedItems"`
	Name              string `xml:"Name,attr"`
}

type Playlist struct {
	Path    string `xml:" Path"`
	Shuffle bool   `xml:" Shuffle"`
	Fader   Fader  `xml:" Fader,omitempty"`
}

type WeekSchedule struct {
	Mon Day `xml:" Mon"`
	Tue Day `xml:" Tue"`
	Wed Day `xml:" Wed"`
	Thu Day `xml:" Thu"`
	Fri Day `xml:" Fri"`
	Sat Day `xml:" Sat"`
	Sun Day `xml:" Sun"`
}

type Zone struct {
	Maintainer   string                 `xml:" Maintainer,omitempty"`
	Description  string                 `xml:" Description,omitempty"`
	Comment      string                 `xml:" Comment,omitempty"`
	Main         Playlist               `xml:" Main"`
	Fallback     Playlist               `xml:" Fallback,omitempty"`
	Intermediate []IntermediatePlaylist `xml:" Intermediate,omitempty"`
	Name         string                 `xml:"Name,attr"`
	Start        time.Time              `xml:"Start,attr"`
}

func (t *Zone) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T Zone
	var layout struct {
		*T
		Start *xsdTime `xml:"Start,attr"`
	}
	layout.T = (*T)(t)
	layout.Start = (*xsdTime)(&layout.T.Start)
	return e.EncodeElement(layout, start)
}
func (t *Zone) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Zone
	var overlay struct {
		*T
		Start *xsdTime `xml:"Start,attr"`
	}
	overlay.T = (*T)(t)
	overlay.Start = (*xsdTime)(&overlay.T.Start)
	return d.DecodeElement(&overlay, &start)
}

type xsdTime time.Time

func (t *xsdTime) UnmarshalText(text []byte) error {
	return _unmarshalTime(text, (*time.Time)(t), "15:04:05.999999999")
}
func (t xsdTime) MarshalText() ([]byte, error) {
	return []byte((time.Time)(t).Format("15:04:05.999999999")), nil
}
func (t xsdTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if (time.Time)(t).IsZero() {
		return nil
	}
	m, err := t.MarshalText()
	if err != nil {
		return err
	}
	return e.EncodeElement(m, start)
}
func (t xsdTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if (time.Time)(t).IsZero() {
		return xml.Attr{}, nil
	}
	m, err := t.MarshalText()
	return xml.Attr{Name: name, Value: string(m)}, err
}
func _unmarshalTime(text []byte, t *time.Time, format string) (err error) {
	s := string(bytes.TrimSpace(text))
	*t, err = time.Parse(format, s)
	if _, ok := err.(*time.ParseError); ok {
		*t, err = time.Parse(format+"Z07:00", s)
	}
	return err
}
