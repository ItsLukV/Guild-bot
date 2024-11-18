package guildData

import (
	"fmt"
	"time"
)

// -----------------------------------------------
// --------------- Main event type ---------------
// -----------------------------------------------

type Event interface {
	GetId() int
	GetType() EventType
	GetDescription() string
	GetDuration() int
	GetIsActive() bool
	GetStartTime() time.Time
	Start() error
	End() error
	AddUser(user GuildUser) error
}

type EventType int

const (
	// Slayer   EventType = "slayer"
	// Diana    EventType = "diana"
	// Dungeons EventType = "dungeons"
	Slayer   EventType = 0
	Diana    EventType = 1
	Dungeons EventType = 2
)

func (e EventType) String() string {
	switch e {
	case Slayer:
		return "slayer"
	case Diana:
		return "diana"
	case Dungeons:
		return "dungeons"
	default:
		return fmt.Sprintf("%d", e)
	}
}

type GuildEvent struct {
	Id            int
	Type          EventType
	Description   string
	StartTime     time.Time
	DurationHours int
	IsActive      bool
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

func (g *GuildEvent) GetIsActive() bool {
	return g.IsActive
}

func (g *GuildEvent) GetStartTime() time.Time {
	return g.StartTime
}

func (g *GuildEvent) GetDuration() int {
	return g.DurationHours
}

func (g *GuildEvent) Start() error {
	// Common start logic
	return nil
}

func (g *GuildEvent) End() error {
	// Common start logic
	return nil
}

func (g *GuildEvent) AddUser(user GuildUser) error {
	// Common add user logic
	return nil
}

// -----------------------------------------------
// ---------------- Costum events ----------------
// -----------------------------------------------

type SlayerEventData struct {
	User             string
	Date             time.Time
	ZombieTier1Kills int
	ZombieTier2Kills int
}

type SlayerEvent struct {
	GuildEvent
	Data map[string]SlayerEventData
}

func (g *SlayerEvent) AddUser(user GuildUser) error {
	// Common add user logic
	return nil
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

func NewGuildEvent(id int, eventType EventType, description string, startTime time.Time, durationHours int) Event {
	baseEvent := GuildEvent{
		Id:            id,
		Type:          eventType,
		Description:   description,
		StartTime:     startTime,
		DurationHours: durationHours,
		IsActive:      false,
	}

	switch eventType {
	case Slayer:
		return &SlayerEvent{
			GuildEvent: baseEvent,
			Data:       make(map[string]SlayerEventData),
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

// Example data
// "slayer_bosses": {
// 	"zombie": {
// 		"claimed_levels": {
// 			"level_1": true,
// 			"level_2": true,
// 			"level_3": true,
// 			"level_4": true,
// 			"level_5": true,
// 			"level_6": true,
// 			"level_7_special": true,
// 			"level_8_special": true,
// 			"level_9_special": true
// 		},
// 		"boss_kills_tier_0": 128,
// 		"xp": 3896465,
// 		"boss_kills_tier_1": 185,
// 		"boss_kills_tier_2": 32,
// 		"boss_kills_tier_3": 293,
// 		"boss_kills_tier_4": 2479,
// 		"boss_attempts_tier_4": 72,
// 		"boss_attempts_tier_0": 39
// 	},
// 	"spider": {
// 		"claimed_levels": {
// 			"level_1": true,
// 			"level_2": true,
// 			"level_3": true,
// 			"level_4": true,
// 			"level_5": true,
// 			"level_6": true,
// 			"level_7": true,
// 			"level_8": true,
// 			"level_9": true
// 		},
// 		"boss_kills_tier_0": 10,
// 		"xp": 1000000,
// 		"boss_kills_tier_1": 36,
// 		"boss_kills_tier_2": 28,
// 		"boss_kills_tier_3": 1982,
// 		"boss_attempts_tier_3": 1940,
// 		"boss_attempts_tier_2": 5,
// 		"boss_attempts_tier_1": 14
// 	},
