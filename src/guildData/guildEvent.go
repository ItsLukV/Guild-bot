package guildData

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// -----------------------------------------------
// --------------- Main event type ---------------
// -----------------------------------------------

type EventSaver interface {
	SaveEventData(event Event) error
}

type Event interface {
	String() string
	GetId() string
	GetEventName() string
	GetType() EventType
	GetDescription() string
	GetIsActive() bool
	GetDuration() int
	GetStartTime() time.Time
	GetUsers() map[Snowflake]*GuildUser
	GetLastFetch() time.Time
	SetLastFetch(time time.Time)
	Start() error
	End() error
	AddUser(user *GuildUser) error
	ShouldFetchData() bool
	IsHidden() bool
	sortByPoints() ([]Snowflake, []float64)
	FetchData() error
	HasEnded() bool
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
	EventName     string
	Type          EventType
	Description   string
	StartTime     time.Time
	DurationHours int
	IsActive      bool
	LastFetch     time.Time
	Users         map[Snowflake]*GuildUser
	Hidden        bool
	hasEnded      bool
}

func (g *GuildEvent) String() string {
	return fmt.Sprintf(
		"{ Id: %v,\n EventName: %v,\n Type: %v,\n Description %v,\n StartTime %v,\n DurationHours %v,\n IsActive: %v,\n"+
			" LastFetch %v,\n Users: %v,\n Hidden: %v \n}",
		g.Id,
		g.EventName,
		g.Type,
		g.Description,
		g.StartTime,
		g.DurationHours,
		g.IsActive,
		g.LastFetch,
		g.Users,
		g.Hidden,
	)
}

func (g *GuildEvent) GetId() string {
	return g.Id
}

func (g *GuildEvent) GetEventName() string {
	return g.EventName
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

func (g *GuildEvent) GetUsers() map[Snowflake]*GuildUser {
	return g.Users
}

func (g *GuildEvent) GetLastFetch() time.Time {
	return g.LastFetch
}

func (g *GuildEvent) SetLastFetch(time time.Time) {
	g.LastFetch = time
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

func (g *GuildEvent) ShouldFetchData() bool {
	now := time.Now().UTC()
	elapsed := now.Sub(g.LastFetch)
	// fmt.Printf("Current Time: %v\n", now)
	// fmt.Printf("Last Fetch: %v\n", g.LastFetch)
	// fmt.Printf("Elapsed Time: %v\n", elapsed)

	if elapsed < 0 {
		log.Println("Warning: LastFetch is in the future.")
		return false
	}

	return elapsed > 3*time.Hour
}

func (g *GuildEvent) IsHidden() bool {
	return g.Hidden
}

func GetLeaderboard(event Event) *discordgo.MessageEmbed {
	snowflakes, points := event.sortByPoints()
	// Build the leaderboard embed
	var builder strings.Builder
	for i, snowflake := range snowflakes {
		builder.WriteString(fmt.Sprintf("`%v.` <@%v>: %.2f\n", i+1, snowflake, points[i]))
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s Event Leaderboard", event.GetEventName()),
		Description: fmt.Sprintf("Here are the results from the %s event.", event.GetType().String()),
		Color:       0x3242a8,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Leaderboard",
				Value: builder.String(),
			},
		},
	}
	return embed
}

func (g *GuildEvent) sortByPoints() ([]Snowflake, []float64) {
	log.Panicln("This guild event doesn't implement sortByPoints")
	return nil, nil
}

func (s *GuildEvent) FetchData() error {
	return nil
}

func (s *GuildEvent) HasEnded() bool {
	return s.hasEnded
}

// -----------------------------------------------
// ---------------- Custom events ----------------
// -----------------------------------------------

type SlayerEventData struct {
	Id        string
	FetchDate time.Time
	BossData  map[BossType]SlayerBossData // Per boss data
}

func (s *SlayerEventData) calculatePoints() float64 {
	totalPoints := float64(0)
	slayerXp := float64(0)
	for bossType, bossData := range s.BossData {
		bossModifier := float64(1)
		xpModifier := float64(1)
		points := float64(0)
		switch bossType {

		case Blaze:
			bossModifier = 935.0455
			xpModifier = 10
		case Enderman:
			bossModifier = 996.3003
			xpModifier = 10
		case Spider:
			bossModifier = 7019.57
			xpModifier = 3.6
		case Vampire:
			bossModifier = 935.0455
			xpModifier = 10
		case Wolf:
			bossModifier = 2982.06
			xpModifier = 10
		case Zombie:
			bossModifier = 9250
			xpModifier = 1.6
		default:
			log.Printf("Unexpected boss type: %s", bossType.String())
		}
		// Calculate points for kills
		points += float64(bossData.BossKillsTier0)
		points += float64(bossData.BossKillsTier1)
		points += float64(bossData.BossKillsTier2)
		points += float64(bossData.BossKillsTier3)
		points += float64(bossData.BossKillsTier4)

		points /= bossModifier * 5

		// Calculate points for XP
		slayerXp += float64(bossData.Xp) * xpModifier

		totalPoints += points
	}
	totalPoints += slayerXp / 1_000_000
	return totalPoints * 100
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
	Data map[Snowflake][]SlayerEventData // Per player data
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
		slayerEvent.Data[snowflake] = append(slayerEvent.Data[snowflake], SlayerEventData{
			Id:        id,
			FetchDate: time.Now().UTC(),
			BossData:  slayerData,
		},
		)
	}
	slayerEvent.IsActive = true
	slayerEvent.hasEnded = true
	return nil
}

func (slayerEvent *SlayerEvent) End() error {
	log.Printf("Ended a event with %v players", len(slayerEvent.Users))
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
		slayerEvent.Data[snowflake] = append(slayerEvent.Data[snowflake], SlayerEventData{
			Id:        id,
			FetchDate: time.Now().UTC(),
			BossData:  slayerData,
		},
		)
	}
	slayerEvent.IsActive = false
	return nil
}

func (g *SlayerEvent) sortByPoints() ([]Snowflake, []float64) {
	// Create a slice to hold users and their points
	type userPoints struct {
		userId Snowflake
		points float64
	}

	userPointsList := make([]userPoints, 0, len(g.Data))

	// Calculate points for each user and store them
	for user, data := range g.Data {
		oldestPointsAmount := data[0].calculatePoints()
		newestPointsAmount := data[len(data)-1].calculatePoints()

		points := newestPointsAmount - oldestPointsAmount

		userPointsList = append(userPointsList, userPoints{userId: user, points: points})
	}

	// Sort users by points in descending order
	sort.Slice(userPointsList, func(i, j int) bool {
		return userPointsList[i].points > userPointsList[j].points
	})

	// Create separate slices for sorted user IDs and their points
	sortedUsers := make([]Snowflake, len(userPointsList))
	sortedPoints := make([]float64, len(userPointsList))
	for i, userPoint := range userPointsList {
		sortedUsers[i] = userPoint.userId
		sortedPoints[i] = userPoint.points
	}

	return sortedUsers, sortedPoints
}

func (s *SlayerEvent) FetchData() error {
	for _, user := range s.Users {
		mcUuid := user.McUUID
		activeProfile, err := utils.FetchActivePlayerProfile(mcUuid)
		if err != nil {
			return fmt.Errorf("failed to get Active Player Profile: %w", err)
		}

		data, err := utils.FetchPlayerSlayerData(mcUuid, activeProfile)
		if err != nil {
			return fmt.Errorf("failed to fetch slayer data: %w", err)
		}

		convertedData := convertUtilsSlayerDataToEventSlayerData(data)
		id, err := gonanoid.New()
		if err != nil {
			return fmt.Errorf("failed to create an id: %w", err)
		}

		s.Data[user.Snowflake] = append(s.Data[user.Snowflake], SlayerEventData{
			Id:        id,
			FetchDate: time.Time{},
			BossData:  convertedData,
		})
	}
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

func NewGuildEvent(
	id string,
	eventName string,
	eventType EventType,
	description string,
	startTime time.Time,
	lastFetch time.Time,
	durationHours int,
	isActive bool,
	hidden bool,
	hasEnded bool,
) Event {
	baseEvent := GuildEvent{
		Id:            id,
		EventName:     eventName,
		Type:          eventType,
		Description:   description,
		StartTime:     startTime,
		DurationHours: durationHours,
		IsActive:      isActive,
		LastFetch:     lastFetch,
		Users:         make(map[Snowflake]*GuildUser),
		Hidden:        hidden,
		hasEnded:      hasEnded,
	}

	switch eventType {
	case Slayer:
		playerData := make(map[Snowflake][]SlayerEventData)
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
