package presentation

import (
	"fmt"
	"time"

	"github.com/traPtitech/knoQ/domain"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/ical"
)

type ScheduleStatus int

const (
	Pending ScheduleStatus = iota + 1
	Attendance
	Absent
)

// EventReqWrite is
//go:generate gotypeconverter -s EventReqWrite -d domain.WriteEventParams -o converter.go .
type EventReqWrite struct {
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	AllowTogether bool        `json:"sharedRoom"`
	TimeStart     time.Time   `json:"timeStart"`
	TimeEnd       time.Time   `json:"timeEnd"`
	RoomID        uuid.UUID   `json:"roomId"`
	Place         string      `json:"place"`
	GroupID       uuid.UUID   `json:"groupId"`
	Admins        []uuid.UUID `json:"admins"`
	Tags          []struct {
		Name   string `json:"name"`
		Locked bool   `json:"locked"`
	} `json:"tags"`
	Open bool `json:"open"`
}

type EventTagReq struct {
	Name string `json:"name"`
}

type EventScheduleStatusReq struct {
	Schedule ScheduleStatus `json:"schedule"`
}

// EventDetailRes is experimental
//go:generate gotypeconverter -s domain.Event -d EventDetailRes -o converter.go .
type EventDetailRes struct {
	ID            uuid.UUID          `json:"eventId"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Room          RoomRes            `json:"room"`
	Group         GroupRes           `json:"group"`
	Place         string             `json:"place" cvt:"Room"`
	GroupName     string             `json:"groupName" cvt:"Group"`
	TimeStart     time.Time          `json:"timeStart"`
	TimeEnd       time.Time          `json:"timeEnd"`
	CreatedBy     uuid.UUID          `json:"createdBy"`
	Admins        []uuid.UUID        `json:"admins"`
	Tags          []EventTagRes      `json:"tags"`
	AllowTogether bool               `json:"sharedRoom"`
	Open          bool               `json:"open"`
	Attendees     []EventAttendeeRes `json:"attendees"`
	Model
}

type EventTagRes struct {
	ID     uuid.UUID `json:"tagId" cvt:"Tag"`
	Name   string    `json:"name" cvt:"Tag"`
	Locked bool      `json:"locked"`
}

type EventAttendeeRes struct {
	ID       uuid.UUID      `json:"userId" cvt:"UserID"`
	Schedule ScheduleStatus `json:"schedule"`
}

// EventRes is for multiple response
//go:generate gotypeconverter -s domain.Event -d EventRes -o converter.go .
//go:generate gotypeconverter -s []*domain.Event -d []EventRes -o converter.go .
type EventRes struct {
	ID            uuid.UUID          `json:"eventId"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	AllowTogether bool               `json:"sharedRoom"`
	TimeStart     time.Time          `json:"timeStart"`
	TimeEnd       time.Time          `json:"timeEnd"`
	RoomID        uuid.UUID          `json:"roomId" cvt:"Room"`
	GroupID       uuid.UUID          `json:"groupId" cvt:"Group"`
	Place         string             `json:"place" cvt:"Room"`
	GroupName     string             `json:"groupName" cvt:"Group"`
	Admins        []uuid.UUID        `json:"admins"`
	Tags          []EventTagRes      `json:"tags"`
	CreatedBy     uuid.UUID          `json:"createdBy"`
	Open          bool               `json:"open"`
	Attendees     []EventAttendeeRes `json:"attendees"`
	Model
}

func iCalVeventFormat(e *domain.Event, host string) *ical.Event {
	timeLayout := "20060102T150405Z"
	vevent := ical.NewEvent()
	vevent.AddProperty("uid", e.ID.String())
	vevent.AddProperty("dtstamp", time.Now().UTC().Format(timeLayout))
	vevent.AddProperty("dtstart", e.TimeStart.UTC().Format(timeLayout))
	vevent.AddProperty("dtend", e.TimeEnd.UTC().Format(timeLayout))
	vevent.AddProperty("created", e.CreatedAt.UTC().Format(timeLayout))
	vevent.AddProperty("last-modified", e.UpdatedAt.UTC().Format(timeLayout))
	vevent.AddProperty("summary", e.Name)
	e.Description += "\n\n"
	e.Description += "-----------------------------------\n"
	e.Description += "イベント詳細ページ\n"
	e.Description += fmt.Sprintf("%s/events/%v", host, e.ID)
	vevent.AddProperty("description", e.Description)
	vevent.AddProperty("location", e.Room.Place)
	vevent.AddProperty("organizer", e.CreatedBy.DisplayName)

	return vevent
}

func ICalFormat(events []*domain.Event, host string) *ical.Calendar {
	c := ical.New()
	ical.NewEvent()
	tz := ical.NewTimezone()
	tz.AddProperty("TZID", "Asia/Tokyo")
	std := ical.NewStandard()
	std.AddProperty("TZOFFSETFROM", "+9000")
	std.AddProperty("TZOFFSETTO", "+9000")
	std.AddProperty("TZNAME", "JST")
	std.AddProperty("DTSTART", "19700101T000000")
	tz.AddEntry(std)
	c.AddEntry(tz)

	for _, e := range events {
		vevent := iCalVeventFormat(e, host)
		c.AddEntry(vevent)
	}
	return c
}
