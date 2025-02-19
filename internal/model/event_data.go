package model

import (
	"fmt"
	"sort"
	"time"
)

type EventData interface {
	GetUser() string
	FormatToDiscordDisplay() string
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

func (d DianaData) FormatToDiscordDisplay() string {
	userName, err := FetchUsername(d.Id)
	if err != nil {
		userName = d.Id
	}

	// Example using Discord markdown + emojis to highlight important stats
	return fmt.Sprintf(
		`**%s** *(Diana)*
		> **Fetched**: %s
		> :coin: **BurrowsTreasure**:  %d
		> :crossed_swords: **BurrowsCombat**:   %d
		> :sparkles: **GaiaConstruct**:   %d
		> :dragon: **MinosChampion**:     %d
		> :bow_and_arrow: **MinosHunter**:       %d
		> :dragon_face: **MinosInquisitor**:    %d
		> :cow2: **Minotaur**:           %d
		> :cat2: **SiameseLynx**:        %d

		`,
		userName,
		d.FetchTime.Format(time.RFC1123),
		d.BurrowsTreasure,
		d.BurrowsCombat,
		d.GaiaConstruct,
		d.MinosChampion,
		d.MinosHunter,
		d.MinosInquisitor,
		d.Minotaur,
		d.SiameseLynx,
	)
}

type DungeonsData struct {
	ID                string             `json:"id"`
	FetchTime         time.Time          `json:"fetch_time"`
	Experience        float64            `json:"experience"`
	Completions       map[int]int        `json:"completions"`
	MasterCompletions map[int]int        `json:"master_completions"`
	ClassXP           map[string]float64 `json:"class_xp"`
	Secrets           int                `json:"secrets"`
}

func (d DungeonsData) GetUser() string {
	return d.ID
}

func (d DungeonsData) FormatToDiscordDisplay() string {
	userName, err := FetchUsername(d.ID)
	if err != nil {
		userName = d.ID
	}

	// We’ll build a string for Completions, MasterCompletions, and ClassXP
	completionsText := formatCompletions(d.Completions)
	masterCompletionsText := formatCompletions(d.MasterCompletions)
	classXPText := formatClassXP(d.ClassXP)

	// Use Discord Markdown for a neat output
	return fmt.Sprintf(
		`**%s** *(Dungeons)*
		> **Fetched**: %s
		> :star: **Experience**:  %.2f
		> :key:  **Secrets**:     %d

		> :crossed_swords: **Completions**:
		%s

		> :dragon_face: **MasterCompletions**:
		%s

		> :book: **ClassXP**:
		%s

		`,
		userName,
		d.FetchTime.Format(time.RFC1123),
		d.Experience,
		d.Secrets,
		completionsText,
		masterCompletionsText,
		classXPText,
	)
}

func formatCompletions(completions map[int]int) string {
	if len(completions) == 0 {
		return ">    *None*"
	}

	// We can sort floors numerically for consistent output
	var floors []int
	for floor := range completions {
		floors = append(floors, floor)
	}
	sort.Ints(floors)

	// Build lines
	result := ""
	for _, floor := range floors {
		count := completions[floor]
		// Each line is further quoted with "> " for a nice Discord blockquote effect
		result += fmt.Sprintf(">    • Floor %d: %d\n", floor, count)
	}
	return result
}

func formatClassXP(classXP map[string]float64) string {
	if len(classXP) == 0 {
		return ">    *None*"
	}

	// Convert XP to levels
	classLevels := toClassLvl(classXP)

	// If you want alphabetical sorting of class names:
	var classes []string
	for c := range classLevels {
		classes = append(classes, c)
	}
	sort.Strings(classes)

	result := ""
	for _, className := range classes {
		level := classLevels[className]
		result += fmt.Sprintf(">    • %s: %.2f\n", className, level)
	}
	return result
}

func toClassLvl(classXP map[string]float64) map[string]float64 {
	// XP needed to go from level i to level i+1
	xpRequired := []int{
		50, 75, 110, 160, 230, 330, 470, 670, 950, 1340, 1890, 2665, 3760, 5260, 7380,
		10300, 14400, 20000, 27600, 38000, 52500, 71500, 97000, 132000, 180000, 243000,
		328000, 445000, 600000, 800000, 1065000, 1410000, 1900000, 2500000, 3300000,
		4300000, 5600000, 7200000, 9200000, 12000000, 15000000, 19000000, 24000000,
		30000000, 38000000, 48000000, 60000000, 75000000, 93000000, 116250000,
	}

	// We'll return a map from className -> fractionalLevel
	fractionalLevels := make(map[string]float64)

	for class, xp := range classXP {
		var lvl float64 // The final fractional level
		var totalSoFar float64

		// Walk through each threshold
		for i, needed := range xpRequired {
			neededF := float64(needed)

			// Check if the user's XP is below (totalSoFar + neededF)
			if xp < totalSoFar+neededF {
				// Then we are partway through this level
				fraction := (xp - totalSoFar) / neededF // how far we are between i and i+1
				lvl = float64(i) + fraction             // e.g. i=0 => lvl=0.98 for xp=49
				break
			}

			// Otherwise, they've fully cleared this level, add to totalSoFar
			totalSoFar += neededF

			// If they're past the highest threshold, they're effectively at max
			if i == len(xpRequired)-1 && xp >= totalSoFar {
				lvl = float64(len(xpRequired))
			}
		}

		// If the user doesn't even exceed the first threshold,
		// lvl remains 0 or the partial fraction assigned above.
		fractionalLevels[class] = lvl
	}

	return fractionalLevels
}
