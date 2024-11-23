CREATE TABLE IF NOT EXISTS Users (
    discord_snowflake TEXT PRIMARY KEY,
    discord_username TEXT,
    minecraft_username TEXT,
    minecraft_uuid TEXT
);

CREATE TABLE IF NOT EXISTS GuildEventType (
    id INT PRIMARY KEY,
    type_name TEXT
);

CREATE TABLE IF NOT EXISTS GuildEvent (
    id TEXT PRIMARY KEY,
    event_name TEXT,
    guild_type_id INT REFERENCES GuildEventType(id),
    description TEXT,
    start_time TIMESTAMP,
    last_fetch TIMESTAMP,
    duration_hours INT,
    is_active BOOLEAN,
    hidden BOOLEAN
);

CREATE TABLE IF NOT EXISTS SlayerEventData (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES Users(discord_snowflake),
    fetch_date TIMESTAMP,
    guild_event_id TEXT REFERENCES GuildEvent(id)
);

CREATE TABLE IF NOT EXISTS SlayerBossType (
    id INT PRIMARY KEY,
    slayer_name TEXT
);

CREATE TABLE IF NOT EXISTS BossData (
    id TEXT PRIMARY KEY,
    slayer_event_data_id TEXT REFERENCES SlayerEventData(id),
    slayer_boss_type_id INT REFERENCES SlayerBossType(id),
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
);


INSERT INTO GuildEventType (id, type_name)
VALUES
(0, 'slayer'),
(1, 'diana'),
(2, 'dungeons');


INSERT INTO SlayerBossType (id, slayer_name)
VALUES
(0, 'zombie'),
(1, 'spider'),
(2, 'wolf'),
(3, 'enderman'),
(4, 'vampire'),
(5, 'blaze');
