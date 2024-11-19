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

func LoadGuildBot() (guildData.GuildBot, error) {
	// Fetch users from the database
	users, err := GetInstance().fetchUsers()
	if err != nil {
		return guildData.GuildBot{}, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Fetch events from the database
	guildEvents, err := GetInstance().fetchEvents()
	if err != nil {
		return guildData.GuildBot{}, fmt.Errorf("failed to fetch guild events: %w", err)
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
				return guildData.GuildBot{}, fmt.Errorf("failed to cast event to SlayerEvent")
			}

			// Start a goroutine to fetch and insert slayer data for this event
			go func(e *guildData.SlayerEvent) {
				err := GetInstance().fetchAndInsertSlayerData(e)
				if err != nil {
					log.Printf("Failed to fetch and insert Slayer data for event %d: %v", e.GetId(), err)
				}
			}(slayerEvent)
		default:
			panic("unexpected guildData.EventType")
		}
	}

	// Return the populated GuildBot struct
	return guildData.GuildBot{
		Users:  users,
		Events: guildEvents,
	}, nil
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
        INSERT INTO GuildEvent (id, guild_type_id, description, start_time, duration_hours)
        VALUES ($1, $2, $3, $4, $5)
    `

	_, err := d.pool.Exec(ctx, query,
		event.GetId(),
		event.GetType(),
		event.GetDescription(),
		event.GetStartTime(),
		event.GetDuration(),
	)
	if err != nil {
		log.Printf("Failed to insert guild event %v: %v", event.GetId(), err)
		return fmt.Errorf("failed to insert guild event: %w", err)
	}

	log.Printf("Guild event %v inserted successfully.", event.GetId())
	return nil
}

func (d *Database) SaveStartEventData(event guildData.Event) error {
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

	for userID, eventData := range slayerEvent.Data {
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

	// Insert from temporary table into SlayerEventData (COMMIT needed before SlayerBossData insertion)
	if _, err = tx.Exec(ctx, `
        INSERT INTO SlayerEventData (id, user_id, fetch_date, guild_event_id)
        SELECT id, user_id, fetch_date, guild_event_id FROM temp_slayer_data
		ON CONFLICT (id) DO NOTHING;
    `); err != nil {
		return fmt.Errorf("failed to insert into slayer_event_data: %w", err)
	}

	// Now commit the transaction to ensure that `SlayerEventData` is updated in the database
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit SlayerEventData transaction: %w", err)
	}

	// Start a new transaction for SlayerBossData since it depends on committed SlayerEventData
	tx, err = d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start database transaction for slayer boss data: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

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

	log.Printf("Batch insert of slayer event data completed successfully!")
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
        SELECT id, guild_type_id, description, start_time, duration_hours
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
		var guild_type_id int
		var description string
		var start_time time.Time
		var duration_hours int

		// Scan the row
		err := rows.Scan(&id, &guild_type_id, &description, &start_time, &duration_hours)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Store the user in the map with the snowflake as the key
		events[id] = guildData.NewGuildEvent(id, guildData.EventType(guild_type_id), description, start_time, duration_hours)
		log.Printf("Loaded event: %v, type: %v, description: %v", id, guildData.EventType(guild_type_id), description)
	}

	// Check for any errors that occurred during iteration
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}

	return events, nil
}

func (d *Database) fetchAndInsertSlayerData(slayerEvent *guildData.SlayerEvent) error {
	// Create a context for the query
	ctx := context.Background()

	// Start a transaction
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

	// Fetch slayer event data from the slayer_event_data table
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

	// Iterate through the slayer_event_data rows and populate the SlayerEvent structure
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

		// Fetch associated SlayerEventData for this slayerEventData.Id
		bossQuery := `
            SELECT id, slayer_boss_type_id, boss_kills_tier_0, boss_kills_tier_1, boss_kills_tier_2,
                   boss_kills_tier_3, boss_kills_tier_4, boss_attempts_tier_0, boss_attempts_tier_1,
                   boss_attempts_tier_2, boss_attempts_tier_3, boss_attempts_tier_4, xp
            FROM BossData
            WHERE slayer_event_data_id = $1
        `

		bossRows, err := tx.Query(ctx, bossQuery, slayerEventData.Id)
		if err != nil {
			return fmt.Errorf("failed to query slayer_boss_data for event_data_id %v: %w", slayerEventData.Id, err)
		}
		defer bossRows.Close()

		slayerEventData.BossData = make(map[guildData.BossType]guildData.SlayerBossData)

		for bossRows.Next() {
			var bossData guildData.SlayerBossData
			var bossType guildData.BossType
			err := bossRows.Scan(
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
				return fmt.Errorf("failed to scan slayer_boss_data row: %w", err)
			}

			slayerEventData.BossData[bossType] = bossData
		}

		slayerEvent.Data[userId] = slayerEventData
	}

	// Check for errors that occurred during iteration of the slayer_event_data rows
	if rows.Err() != nil {
		return fmt.Errorf("error occurred during iteration of slayer_event_data rows: %w", rows.Err())
	}

	return nil
}
