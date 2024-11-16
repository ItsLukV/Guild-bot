package guildEvent

import (
	"time"
)

// -----------------------------------------------
// --------------- Main event type ---------------
// -----------------------------------------------

type Event interface {
	GetId() int
	GetType() EventType
	GetDescription() string
	GetDuration() time.Duration
	IsActive() bool
	Start() error
	End() error
	AddUser() error
}

type EventType string

const (
	Slayer   EventType = "slayer"
	Diana    EventType = "diana"
	Dungeons EventType = "dungeons"
)

type GuildEvent struct {
	Id          int
	Type        EventType
	Description string
	// Users
	StartTime time.Time
	EndTime   time.Time
}

func (g *GuildEvent) GetId() int {
	return g.Id
}

func (g *GuildEvent) GetType() EventType {
	return g.Type
}

func (g *GuildEvent) GetDescription() string {
	return g.Description
}

func (g *GuildEvent) GetDuration() time.Duration {
	return g.EndTime.Sub(g.StartTime)
}

func (g *GuildEvent) IsActive() bool {
	now := time.Now()
	return now.After(g.StartTime) && now.Before(g.EndTime)
}

func (g *GuildEvent) Start() error {
	// Common start logic
	return nil
}

func (g *GuildEvent) End() error {
	// Common start logic
	return nil
}

func (g *GuildEvent) AddUser() error {
	// Common add user logic
	return nil
}

// -----------------------------------------------
// ---------------- Costum events ----------------
// -----------------------------------------------

type SlayerEvent struct {
	GuildEvent
}

type DianaEvent struct {
	GuildEvent
}

type DungeonsEvent struct {
	GuildEvent
}

// -----------------------------------------------
// ---------------- Constructors -----------------
// -----------------------------------------------

func NewGuildEvent(eventType EventType, description string, durationHours int) Event {
	startTime := time.Now()
	endTime := startTime.Add(time.Duration(durationHours) * time.Hour)
	baseEvent := GuildEvent{
		Id:          1,
		Type:        eventType,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	switch eventType {
	case Slayer:
		return &SlayerEvent{
			GuildEvent: baseEvent,
		}
	case Diana:
		return &DianaEvent{
			GuildEvent: baseEvent,
		}
	case Dungeons:
		return &DungeonsEvent{
			GuildEvent: baseEvent,
		}
	default:
		return &baseEvent
	}
}
