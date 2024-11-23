package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// -----------------------------------------------
// -------------- Struct and const --------------
// -----------------------------------------------

var (
	host            = os.Getenv("DB_HOST")
	port     uint16 = 5432
	user            = os.Getenv("DB_USER")
	password        = os.Getenv("DB_PASSWORD")
	dbname          = os.Getenv("DB")
)

type Database struct {
	pool *pgxpool.Pool
}

// -----------------------------------------------
// ------------------ Functions ------------------
// -----------------------------------------------

var lock = &sync.Mutex{}

var singleInstance *Database

func GetInstance() *Database {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("Creating single instance now.")
			singleInstance = &Database{}
			singleInstance.init()
		} else {
			log.Println("Single instance already created.")
		}
	} else {
		// log.Println("Single instance already created.")
	}

	return singleInstance
}

func (d *Database) init() {
	if d == nil {
		return
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	dbpool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	} else {
		log.Println("Successfully connected!")
	}

	d.pool = dbpool

}

func (d *Database) Close() {
	d.pool.Close()
}

func LoadGuildBot(guildBot *guildData.GuildBot) error {
	// Fetch users from the database
	users, err := GetInstance().fetchUsers()
	if err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	// Fetch events from the database
	guildEvents, err := GetInstance().fetchEvents()
	if err != nil {
		return fmt.Errorf("failed to fetch guild events: %w", err)
	}

	// Iterate over fetched guild events
	for i := range guildEvents {
		event := guildEvents[i]

		switch event.GetType() {
		case guildData.Diana:
			// Handling Diana event
			// You can add specific processing for Diana here if needed

		case guildData.Dungeons:
			// Handling Dungeons event
			// You can add specific processing for Dungeons here if needed

		case guildData.Slayer:
			// Type assertion to ensure this event is treated as a SlayerEvent
			slayerEvent, ok := event.(*guildData.SlayerEvent)
			if !ok {
				return fmt.Errorf("failed to cast event to SlayerEvent")
			}

			// Start a goroutine to fetch and insert slayer data for this event
			go func(e *guildData.SlayerEvent) {
				err := GetInstance().fetchSlayerData(e)
				if err != nil {
					log.Printf("Failed to fetch Slayer data for event %s: %v", e.GetId(), err)
				}
			}(slayerEvent)
		default:
			panic("unexpected guildData.EventType")
		}
	}

	// Return the populated GuildBot struct
	guildBot.Users = users
	guildBot.Events = guildEvents
	guildBot.EventSaver = GetInstance()
	return nil
}

// -----------------------------------------------
// ------------- Insertion functions -------------
// -----------------------------------------------

func (d *Database) SaveUser(user guildData.GuildUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout
	defer cancel()                                                          // Ensure the context is cancelled

	query := `
        INSERT INTO Users (discord_snowflake, discord_username, minecraft_username, minecraft_uuid)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (discord_snowflake) DO NOTHING;
    `

	_, err := d.pool.Exec(ctx, query,
		user.Snowflake,       // $1: discord_snowflake
		user.DiscordUsername, // $2: discord_username
		user.McUsername,      // $3: minecraft_username
		user.McUUID,          // $4: minecraft_uuid
	)
	if err != nil {
		log.Printf("Failed to insert user %s: %v", user.Snowflake, err)
		return fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("User %s inserted successfully.", user.Snowflake)
	return nil
}

// Saves the event to the db
func (d *Database) SaveEvent(event guildData.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout
	defer cancel()                                                          // Ensure the context is cancelled

	query := `
        INSERT INTO GuildEvent (id, event_name, guild_type_id, description, start_time, last_fetch, duration_hours, is_active, hidden, has_ended)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := d.pool.Exec(ctx, query,
		event.GetId(),
		event.GetEventName(),
		event.GetType(),
		event.GetDescription(),
		event.GetStartTime(),
		event.GetLastFetch(),
		event.GetDuration(),
		event.GetIsActive(),
		event.IsHidden(),
		event.HasEnded(),
	)
	if err != nil {
		log.Printf("Failed to insert guild event %v: %v", event.GetId(), err)
		return fmt.Errorf("failed to insert guild event: %w", err)
	}

	log.Printf("Guild event %v inserted successfully.", event.GetId())
	return nil
}

// UpdateEvent updates the existing event in the database
func (d *Database) UpdateEvent(event guildData.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Set a timeout
	defer cancel()                                                          // Ensure the context is cancelled

	query := `
        UPDATE GuildEvent
        SET event_name = $2,
            guild_type_id = $3,
            description = $4,
            start_time = $5,
            last_fetch = $6,
            duration_hours = $7,
            is_active = $8,
            hidden = $9
			has_ended = $10
        WHERE id = $1
    `

	// Execute the update query with the appropriate values from the event
	_, err := d.pool.Exec(ctx, query,
		event.GetId(),          // $1: ID of the event to update
		event.GetEventName(),   // $2: Event name
		event.GetType(),        // $3: Type of the event
		event.GetDescription(), // $4: Description of the event
		event.GetStartTime(),   // $5: Start time of the event
		event.GetLastFetch(),   // $6: Last fetch time
		event.GetDuration(),    // $7: Duration in hours
		event.GetIsActive(),    // $8: Is the event active?
		event.IsHidden(),       // $9: Is the event hidden?
		event.HasEnded(),       // $10: has the event ended
	)
	if err != nil {
		log.Printf("Failed to update guild event %v: %v", event.GetId(), err)
		return fmt.Errorf("failed to update guild event: %w", err)
	}

	log.Printf("Guild event %v updated successfully.", event.GetId())
	return nil
}

func (d *Database) SaveEventData(event guildData.Event) error {
	var wg sync.WaitGroup // Create a WaitGroup to manage goroutines

	switch event.GetType() {
	case guildData.Diana:
		// Handling Diana event
		// You can add specific processing for Diana here if needed

	case guildData.Dungeons:
		// Handling Dungeons event
		// You can add specific processing for Dungeons here if needed

	case guildData.Slayer:
		// Type assertion to ensure this event is treated as a SlayerEvent
		slayerEvent, ok := event.(*guildData.SlayerEvent)
		if !ok {
			return fmt.Errorf("failed to cast event to SlayerEvent")
		}

		// Increment WaitGroup counter
		wg.Add(1)

		// Start a goroutine to fetch and insert slayer data for this event
		go func(e *guildData.SlayerEvent) {
			defer wg.Done() // Mark this goroutine as done when finished

			err := GetInstance().SaveSlayerEventData(e)
			if err != nil {
				log.Printf("Failed to insert Slayer data for event %v: %v", e.GetId(), err)
			}
			log.Printf("Inserted slayer data for %v", e.Id)
		}(slayerEvent)

	default:
		panic("unexpected guildData.EventType")
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return nil
}

func (d *Database) SaveSlayerEventData(slayerEvent *guildData.SlayerEvent) error {
	ctx := context.Background()

	// Begin a transaction
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			log.Printf("Batch insert of slayer event data completed successfully!")
		}
	}()

	// Create temporary tables
	if _, err = tx.Exec(ctx, `
        CREATE TEMP TABLE temp_slayer_data (
            id TEXT,
            user_id TEXT,
            fetch_date TIMESTAMP,
            guild_event_id TEXT
        ) ON COMMIT DROP;
    `); err != nil {
		return fmt.Errorf("failed to create temp_slayer_data table: %w", err)
	}

	// Prepare data for bulk insertion
	slayerDataRows := make([][]interface{}, 0, len(slayerEvent.Data))
	slayerBossDataRows := make([][]interface{}, 0)

	for userID, eventDataList := range slayerEvent.Data {
		for _, eventData := range eventDataList {
			// Prepare data for temp_slayer_data
			slayerDataRows = append(slayerDataRows, []interface{}{
				eventData.Id,        // id
				userID,              // user_id
				eventData.FetchDate, // fetch_date
				slayerEvent.Id,      // guild_event_id
			})

			// Prepare data for temp_slayer_boss_data
			for bossType, bossData := range eventData.BossData {
				slayerBossDataRows = append(slayerBossDataRows, []interface{}{
					bossData.Id,                // id
					eventData.Id,               // slayer_event_data_id
					bossType,                   // slayer_boss_type
					bossData.BossKillsTier0,    // tier_0_kills
					bossData.BossKillsTier1,    // tier_1_kills
					bossData.BossKillsTier2,    // tier_2_kills
					bossData.BossKillsTier3,    // tier_3_kills
					bossData.BossKillsTier4,    // tier_4_kills
					bossData.BossAttemptsTier0, // tier_0_attempts
					bossData.BossAttemptsTier1, // tier_1_attempts
					bossData.BossAttemptsTier2, // tier_2_attempts
					bossData.BossAttemptsTier3, // tier_3_attempts
					bossData.BossAttemptsTier4, // tier_4_attempts
					bossData.Xp,                // xp
				})
			}
		}
	}

	// Bulk insert into temp_slayer_data
	columnsSlayerData := []string{"id", "user_id", "fetch_date", "guild_event_id"}
	if _, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"temp_slayer_data"},
		columnsSlayerData,
		pgx.CopyFromRows(slayerDataRows),
	); err != nil {
		return fmt.Errorf("failed to copy data into temp_slayer_data: %w", err)
	}

	// Insert from temporary table into SlayerEventData
	if _, err = tx.Exec(ctx, `
        INSERT INTO SlayerEventData (id, user_id, fetch_date, guild_event_id)
        SELECT id, user_id, fetch_date, guild_event_id FROM temp_slayer_data
		ON CONFLICT (id) DO NOTHING;
    `); err != nil {
		return fmt.Errorf("failed to insert into slayer_event_data: %w", err)
	}

	if _, err = tx.Exec(ctx, `
        CREATE TEMP TABLE temp_slayer_boss_data (
			id TEXT,
            slayer_event_data_id TEXT,
            slayer_boss_type_id INT,
            boss_kills_tier_0 INT NOT NULL DEFAULT 0,
            boss_kills_tier_1 INT NOT NULL DEFAULT 0,
            boss_kills_tier_2 INT NOT NULL DEFAULT 0,
            boss_kills_tier_3 INT NOT NULL DEFAULT 0,
            boss_kills_tier_4 INT NOT NULL DEFAULT 0,
            boss_attempts_tier_0 INT NOT NULL DEFAULT 0,
            boss_attempts_tier_1 INT NOT NULL DEFAULT 0,
            boss_attempts_tier_2 INT NOT NULL DEFAULT 0,
            boss_attempts_tier_3 INT NOT NULL DEFAULT 0,
            boss_attempts_tier_4 INT NOT NULL DEFAULT 0,
            xp INT NOT NULL DEFAULT 0
        ) ON COMMIT DROP;
    `); err != nil {
		return fmt.Errorf("failed to create temp_slayer_boss_data table: %w", err)
	}

	// Bulk insert into temp_slayer_boss_data
	columnsSlayerBossData := []string{
		"id",
		"slayer_event_data_id",
		"slayer_boss_type_id",
		"boss_kills_tier_0",
		"boss_kills_tier_1",
		"boss_kills_tier_2",
		"boss_kills_tier_3",
		"boss_kills_tier_4",
		"boss_attempts_tier_0",
		"boss_attempts_tier_1",
		"boss_attempts_tier_2",
		"boss_attempts_tier_3",
		"boss_attempts_tier_4",
		"xp",
	}
	if _, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"temp_slayer_boss_data"},
		columnsSlayerBossData,
		pgx.CopyFromRows(slayerBossDataRows),
	); err != nil {
		return fmt.Errorf("failed to copy data into temp_slayer_boss_data: %w", err)
	}

	// Insert from temporary table into SlayerBossData
	if _, err = tx.Exec(ctx, `
        INSERT INTO BossData (
            id,
            slayer_event_data_id,
            slayer_boss_type_id,
            boss_kills_tier_0,
            boss_kills_tier_1,
            boss_kills_tier_2,
            boss_kills_tier_3,
            boss_kills_tier_4,
            boss_attempts_tier_0,
            boss_attempts_tier_1,
            boss_attempts_tier_2,
            boss_attempts_tier_3,
            boss_attempts_tier_4,
            xp
        )
        SELECT
            id,
            slayer_event_data_id,
            slayer_boss_type_id,
            boss_kills_tier_0,
            boss_kills_tier_1,
            boss_kills_tier_2,
            boss_kills_tier_3,
            boss_kills_tier_4,
            boss_attempts_tier_0,
            boss_attempts_tier_1,
            boss_attempts_tier_2,
            boss_attempts_tier_3,
            boss_attempts_tier_4,
            xp
        FROM temp_slayer_boss_data
		ON CONFLICT (id) DO NOTHING;
    `); err != nil {
		return fmt.Errorf("failed to insert into slayer_boss_data: %w", err)
	}

	return nil
}

// -----------------------------------------------
// -------------- Fetcher functions --------------
// -----------------------------------------------

func (d *Database) fetchUsers() (map[guildData.Snowflake]guildData.GuildUser, error) {
	// Create a context for the query
	ctx := context.Background()

	// Define the query to fetch all users
	query := `
        SELECT discord_snowflake, discord_username, minecraft_username, minecraft_uuid
        FROM users
    `

	// Execute the query
	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Create a map to store the results
	users := make(map[guildData.Snowflake]guildData.GuildUser)

	// Iterate through the rows and populate the map
	for rows.Next() {
		var user guildData.GuildUser

		// Scan the row into the GuildUser struct and snowflake
		err := rows.Scan(&user.Snowflake, &user.DiscordUsername, &user.McUsername, &user.McUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Store the user in the map with the snowflake as the key
		users[user.Snowflake] = user
	}

	// Check for any errors that occurred during iteration
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}

	return users, nil
}

func (d *Database) fetchEvents() (map[string]guildData.Event, error) {
	// Create a context for the query
	ctx := context.Background()

	// Define the query to fetch all guild events
	query := `
        SELECT id, event_name, guild_type_id, description, start_time, last_fetch, duration_hours, is_active, hidden, has_ended
        FROM GuildEvent
    `

	// Execute the query
	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Create a map to store the results
	events := make(map[string]guildData.Event)

	// Iterate through the rows and populate the map
	for rows.Next() {
		var id string
		var event_name string
		var guild_type_id int
		var description string
		var start_time time.Time
		var last_fetch time.Time
		var duration_hours int
		var is_active bool
		var hidden bool
		var has_ended bool

		// Scan the row
		err := rows.Scan(&id, &event_name, &guild_type_id, &description, &start_time, &last_fetch, &duration_hours, &is_active, &hidden, &has_ended)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Store the user in the map with the snowflake as the key
		events[id] = guildData.NewGuildEvent(id, event_name, guildData.EventType(guild_type_id), description, start_time, last_fetch, duration_hours, is_active, hidden, has_ended)

		log.Println(events[id].String())
	}

	// Check for any errors that occurred during iteration
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}

	return events, nil
}

func (d *Database) fetchSlayerData(slayerEvent *guildData.SlayerEvent) error {
	// Create a context for the query
	ctx := context.Background()

	// Start a transaction for fetching slayer event data
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Fetch slayer event data from the SlayerEventData table
	query := `
        SELECT id, user_id, fetch_date
        FROM SlayerEventData
        WHERE guild_event_id = $1
    `
	rows, err := tx.Query(ctx, query, slayerEvent.Id)
	if err != nil {
		return fmt.Errorf("failed to query slayer_event_data: %w", err)
	}
	defer rows.Close()

	// Create a placeholder for slayer event data and collect user IDs
	slayerEventDataMap := make(map[guildData.Snowflake][]*guildData.SlayerEventData)
	userIdsSet := make(map[guildData.Snowflake]struct{})

	// Iterate through the slayer_event_data rows and populate the SlayerEventData map
	for rows.Next() {
		var slayerEventData guildData.SlayerEventData
		var userId guildData.Snowflake

		err := rows.Scan(
			&slayerEventData.Id,
			&userId,
			&slayerEventData.FetchDate,
		)
		if err != nil {
			return fmt.Errorf("failed to scan slayer_event_data row: %w", err)
		}

		// Append slayerEventData to the user's slice in slayerEventDataMap
		slayerEventDataMap[userId] = append(slayerEventDataMap[userId], &slayerEventData)
		userIdsSet[userId] = struct{}{}
	}

	// Convert userIdsSet to a slice of userIds
	userIds := make([]guildData.Snowflake, 0, len(userIdsSet))
	for userId := range userIdsSet {
		userIds = append(userIds, userId)
	}

	// Check for errors that occurred during iteration of the slayer_event_data rows
	if rows.Err() != nil {
		return fmt.Errorf("error occurred during iteration of slayer_event_data rows: %w", rows.Err())
	}

	// Close the transaction after fetching slayer event data
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction after fetching slayer_event_data: %w", err)
	}

	// Fetch all user data in a single query
	if err := d.fetchSlayerUsers(ctx, userIds, slayerEvent); err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	// Fetch all boss data and associate it with the corresponding slayer event data
	if err := d.fetchBossData(ctx, slayerEventDataMap); err != nil {
		return fmt.Errorf("failed to fetch boss data: %w", err)
	}

	// Add the populated slayerEventData back to the main SlayerEvent
	for userId, slayerEventDataSlice := range slayerEventDataMap {
		// Convert []*SlayerEventData to []SlayerEventData
		slayerEventDataList := make([]guildData.SlayerEventData, len(slayerEventDataSlice))
		for i, slayerEventDataPtr := range slayerEventDataSlice {
			slayerEventDataList[i] = *slayerEventDataPtr
		}
		slayerEvent.Data[userId] = slayerEventDataList
	}

	return nil
}

// fetchSlayerUsers fetches user data for a list of user IDs and adds them to the slayerEvent
func (d *Database) fetchSlayerUsers(ctx context.Context, userIds []guildData.Snowflake, slayerEvent *guildData.SlayerEvent) error {
	if len(userIds) == 0 {
		return nil
	}

	// Start a transaction for fetching user data
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for user data: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	usersQuery := `
        SELECT discord_snowflake, discord_username, minecraft_username, minecraft_uuid
        FROM Users
        WHERE discord_snowflake = ANY($1)
    `

	userRows, err := tx.Query(ctx, usersQuery, userIds)
	if err != nil {
		return fmt.Errorf("failed to query user data: %w", err)
	}
	defer userRows.Close()

	for userRows.Next() {
		var user guildData.GuildUser
		err := userRows.Scan(
			&user.Snowflake,
			&user.DiscordUsername,
			&user.McUsername,
			&user.McUUID,
		)
		if err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}

		slayerEvent.AddUser(&user)
	}

	// Check for errors in the user rows
	if userRows.Err() != nil {
		return fmt.Errorf("error occurred during fetching of users: %w", userRows.Err())
	}

	return nil
}

// fetchBossData fetches boss data and associates it with the corresponding slayer event data
func (d *Database) fetchBossData(ctx context.Context, slayerEventDataMap map[guildData.Snowflake][]*guildData.SlayerEventData) error {
	if len(slayerEventDataMap) == 0 {
		return nil
	}

	// Collect all slayer_event_data IDs and map them to their corresponding SlayerEventData
	slayerEventDataIds := make([]string, 0)
	slayerEventDataIdToData := make(map[string]*guildData.SlayerEventData)
	for _, slayerEventDataSlice := range slayerEventDataMap {
		for _, slayerEventData := range slayerEventDataSlice {
			slayerEventDataIds = append(slayerEventDataIds, slayerEventData.Id)
			slayerEventDataIdToData[slayerEventData.Id] = slayerEventData
		}
	}

	// Start a transaction for fetching boss data
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for boss data: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	bossQuery := `
        SELECT slayer_event_data_id, id, slayer_boss_type_id, boss_kills_tier_0, boss_kills_tier_1, boss_kills_tier_2,
               boss_kills_tier_3, boss_kills_tier_4, boss_attempts_tier_0, boss_attempts_tier_1,
               boss_attempts_tier_2, boss_attempts_tier_3, boss_attempts_tier_4, xp
        FROM BossData
        WHERE slayer_event_data_id = ANY($1)
    `

	bossRows, err := tx.Query(ctx, bossQuery, slayerEventDataIds)
	if err != nil {
		return fmt.Errorf("failed to query boss data: %w", err)
	}
	defer bossRows.Close()

	for bossRows.Next() {
		var bossData guildData.SlayerBossData
		var slayerEventDataId string
		var bossType guildData.BossType
		err := bossRows.Scan(
			&slayerEventDataId,
			&bossData.Id,
			&bossType,
			&bossData.BossKillsTier0,
			&bossData.BossKillsTier1,
			&bossData.BossKillsTier2,
			&bossData.BossKillsTier3,
			&bossData.BossKillsTier4,
			&bossData.BossAttemptsTier0,
			&bossData.BossAttemptsTier1,
			&bossData.BossAttemptsTier2,
			&bossData.BossAttemptsTier3,
			&bossData.BossAttemptsTier4,
			&bossData.Xp,
		)
		if err != nil {
			return fmt.Errorf("failed to scan boss data row: %w", err)
		}

		// Find the corresponding SlayerEventData
		slayerEventData, exists := slayerEventDataIdToData[slayerEventDataId]
		if !exists {
			return fmt.Errorf("slayer event data not found for id %v", slayerEventDataId)
		}

		if slayerEventData.BossData == nil {
			slayerEventData.BossData = make(map[guildData.BossType]guildData.SlayerBossData)
		}

		slayerEventData.BossData[bossType] = bossData
	}

	// Check for errors in the boss rows
	if bossRows.Err() != nil {
		return fmt.Errorf("error occurred during iteration of boss data rows: %w", bossRows.Err())
	}

	return nil
}
