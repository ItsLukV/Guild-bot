package model

import "time"

type EventData interface {
	GetUser() string
}

type DianaData struct {
	Id              string    `json:"id"`
	FetchTime       time.Time `json:"fetch_time"`
	BurrowsTreasure int       `json:"burrows_treasure"`
	BurrowsCombat   int       `json:"burrows_combat"`
	GaiaConstruct   int       `json:"gaia_construct"`
	MinosChampion   int       `json:"minos_champion"`
	MinosHunter     int       `json:"minos_hunter"`
	MinosInquisitor int       `json:"minos_inquisitor"`
	Minotaur        int       `json:"minotaur"`
	SiameseLynx     int       `json:"siamese_lynx"`
}

func (d DianaData) GetUser() string {
	return d.Id
}

type DungeonsData struct {
	ID               string             `json:"id"`
	FetchTime        time.Time          `json:"fetch_time"`
	Experience       float64            `json:"experience"`
	Completions      map[int]int        `json:"completions"`
	MasterCompletions map[int]int       `json:"master_completions"`
	ClassXP          map[string]float64 `json:"class_xp"`
	Secrets          int                `json:"secrets"`
}

func (d DungeonsData) GetUser() string {
	return d.ID
}
