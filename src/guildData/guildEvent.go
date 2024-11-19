package guildData

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ItsLukV/Guild-bot/src/utils"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// -----------------------------------------------
// --------------- Main event type ---------------
// -----------------------------------------------

type Event interface {
	GetId() string
	GetType() EventType
	GetDescription() string
	GetDuration() int
	GetIsActive() bool
	GetStartTime() time.Time
	Start() error
	End() error
	AddUser(user *GuildUser) error
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

type BossType int

const (
	Zombie BossType = iota
	Spider
	Wolf
	Enderman
	Vampire
	Blaze
)

func (e BossType) String() string {
	switch e {
	case Zombie:
		return "zombie"
	case Spider:
		return "spider"
	case Wolf:
		return "wolf"
	case Enderman:
		return "enderman"
	case Vampire:
		return "vampire"
	case Blaze:
		return "blaze"
	default:
		return fmt.Sprintf("%d", e)
	}
}

func ParseBossTypeString(str string) (BossType, bool) {
	bossMap := map[string]BossType{
		"zombie":   Zombie,
		"spider":   Spider,
		"wolf":     Wolf,
		"enderman": Enderman,
		"vampire":  Vampire,
		"blaze":    Blaze,
	}

	boss, ok := bossMap[strings.ToLower(str)]
	return boss, ok
}

type GuildEvent struct {
	Id            string
	Type          EventType
	Description   string
	StartTime     time.Time
	DurationHours int
	IsActive      bool
	Users         map[Snowflake]*GuildUser
}

func (g *GuildEvent) GetId() string {
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
	log.Panicf("This event (%v) has not been implemented", g)
	return nil
}

func (g *GuildEvent) End() error {
	// Common start logic
	return nil
}

func (g *GuildEvent) AddUser(user *GuildUser) error {
	log.Printf("Added user: %v (%v) to %v", user.Snowflake, user.DiscordUsername, g.Id)
	g.Users[user.Snowflake] = user
	return nil
}

// -----------------------------------------------
// ---------------- Custom events ----------------
// -----------------------------------------------

type SlayerEventData struct {
	Id        string
	FetchDate time.Time
	BossData  map[BossType]SlayerBossData // Per boss data
}

type SlayerBossData struct {
	Id                string
	BossKillsTier0    int
	BossKillsTier1    int
	BossKillsTier2    int
	BossKillsTier3    int
	BossKillsTier4    int
	BossAttemptsTier0 int
	BossAttemptsTier1 int
	BossAttemptsTier2 int
	BossAttemptsTier3 int
	BossAttemptsTier4 int
	Xp                int
}

type SlayerEvent struct {
	GuildEvent
	Data map[Snowflake]SlayerEventData // Per player data
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

func NewGuildEvent(id string, eventType EventType, description string, startTime time.Time, durationHours int) Event {
	baseEvent := GuildEvent{
		Id:            id,
		Type:          eventType,
		Description:   description,
		StartTime:     startTime,
		DurationHours: durationHours,
		IsActive:      false,
		Users:         make(map[Snowflake]*GuildUser),
	}

	switch eventType {
	case Slayer:
		playerData := make(map[Snowflake]SlayerEventData)
		slayerEvent := &SlayerEvent{
			GuildEvent: baseEvent,
			Data:       playerData,
		}
		return slayerEvent
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

func (slayerEvent *SlayerEvent) Start() error {
	log.Printf("Started a event with %v players", len(slayerEvent.Users))
	for snowflake, v := range slayerEvent.Users {
		profile, err := utils.FetchActivePlayerProfile(v.McUUID)
		if err != nil {
			return fmt.Errorf("failed to get Active Player Profile: %w", err)
		}

		// Fetch the json data as a utils.SkyblockPlayerSlayerData struct
		v, err := utils.FetchPlayerSlayerData(v.McUUID, profile)
		if err != nil {
			return fmt.Errorf("failed to fetch slayer data: %w", err)
		}

		// Translate the SkyblockPlayerSlayerData struct to
		slayerData := convertUtilsSlayerDataToEventSlayerData(v)
		id, err := gonanoid.New()
		if err != nil {
			log.Printf("Failed to create a Id for SlayerEventData eventId: %v", slayerEvent.Id)
		}
		slayerEvent.Data[snowflake] = SlayerEventData{
			Id:        id,
			FetchDate: time.Now(),
			BossData:  slayerData,
		}
	}
	return nil
}

func convertUtilsSlayerDataToEventSlayerData(slayerData *utils.SkyblockPlayerSlayerData) map[BossType]SlayerBossData {
	out := make(map[BossType]SlayerBossData)
	for k, v := range slayerData.SlayerBosses {
		if bossType, ok := ParseBossTypeString(k); ok {
			id, err := gonanoid.New()
			if err != nil {
				log.Printf("Failed to create a Id for SlayerBossData bossType: %v", bossType)
			}
			out[bossType] = SlayerBossData{
				Id:                id,
				BossKillsTier0:    v.BossKillsTier0,
				BossKillsTier1:    v.BossKillsTier1,
				BossKillsTier2:    v.BossKillsTier2,
				BossKillsTier3:    v.BossKillsTier3,
				BossKillsTier4:    v.BossKillsTier4,
				BossAttemptsTier0: v.BossAttemptsTier0,
				BossAttemptsTier1: v.BossAttemptsTier1,
				BossAttemptsTier2: v.BossAttemptsTier2,
				BossAttemptsTier3: v.BossAttemptsTier3,
				BossAttemptsTier4: v.BossAttemptsTier4,
				Xp:                v.Xp,
			}
		} else {
			// Handle invalid BossType strings if necessary
			fmt.Printf("Invalid BossType: %s\n", k)
		}
	}
	return out
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
