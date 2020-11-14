package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Event struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Room          Room
	Group         Group
	TimeStart     time.Time
	TimeEnd       time.Time
	CreatedBy     User
	Tags          []EventTag
	AllowTogether bool
	Model
}

type EventTag struct {
	Tag
	Locked bool
}

// for repository

// ConInfo is Connection infomation
type ConInfo struct {
	Token     string
	ReqUserID uuid.UUID
}

// WriteEventParams is used create and update
type WriteEventParams struct {
	Name          string
	Description   string
	GroupID       uuid.UUID
	RoomID        uuid.UUID
	TimeStart     time.Time
	TimeEnd       time.Time
	AllowTogether bool
	Tags          []struct {
		Name   string
		Locked bool
	}
}

// WriteTagRelationParams is used create and update
type WriteTagRelationParams struct {
	ID     uuid.UUID
	Locked bool
}

// EventRepository is implemented by ...
type EventRepository interface {
	CreateEvent(eventParams WriteEventParams, info *ConInfo) (*Event, error)

	UpdateEvent(eventID uuid.UUID, eventParams WriteEventParams) (*Event, error)
	AddTagToEvent(eventID uuid.UUID, tagID uuid.UUID, locked bool) error

	DeleteEvent(eventID uuid.UUID) error
	// DeleteTagInEvent delete a tag in that Event
	DeleteTagInEvent(eventID uuid.UUID, tagID uuid.UUID, deleteLocked bool) error

	GetEvent(eventID uuid.UUID) (*Event, error)

	// TODO 一つにまとめる
	GetAllEvents(start *time.Time, end *time.Time) ([]*Event, error)
	GetEventsByGroupIDs(groupIDs []uuid.UUID) ([]*Event, error)
	GetEventsByRoomIDs(roomIDs []uuid.UUID) ([]*Event, error)

	GetEventActivities(day int) ([]*Event, error)
	// GetEventsByFilter allows you to filter the events under any condition.
	// However, you may get an error at runtime.
	GetEventsByFilter(query string, args []interface{}) ([]*Event, error)
}
