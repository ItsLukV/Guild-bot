CREATE TABLE IF NOT EXISTS Users (
    discord_snowflake TEXT PRIMARY KEY,
    discord_username TEXT,
    minecraft_username TEXT,
    minecraft_uuid TEXT
);

CREATE TABLE IF NOT EXISTS GuildEventType (
    id SERIAL PRIMARY KEY,
    type_name TEXT
);

CREATE TABLE IF NOT EXISTS GuildEvent (
    id SERIAL PRIMARY KEY,
    guild_type_id INT REFERENCES GuildEventType(id),
    description TEXT,
    start_time TIMESTAMP,
    duration_hours INT
);

CREATE TABLE IF NOT EXISTS SlayerEventData (
    id SERIAL PRIMARY KEY,
    user_id TEXT REFERENCES Users(discord_snowflake),
    fetch_date TIMESTAMP,
    guild_event INT REFERENCES GuildEvent(id)
);

CREATE TABLE IF NOT EXISTS SlayerBossType (
    id SERIAL PRIMARY KEY,
    slayer_name TEXT
);

CREATE TABLE IF NOT EXISTS BossData (
    id SERIAL PRIMARY KEY,
    slayer_event_id INT REFERENCES SlayerEventData(id),
    slayer_boss_type INT REFERENCES SlayerBossType(id),
    tier_0_kills INT NOT NULL DEFAULT 0,
    tier_1_kills INT NOT NULL DEFAULT 0,
    tier_2_kills INT NOT NULL DEFAULT 0,
    tier_3_kills INT NOT NULL DEFAULT 0,
    tier_4_kills INT NOT NULL DEFAULT 0,
    tier_0_attempts INT NOT NULL DEFAULT 0,
    tier_1_attempts INT NOT NULL DEFAULT 0,
    tier_2_attempts INT NOT NULL DEFAULT 0,
    tier_3_attempts INT NOT NULL DEFAULT 0,
    tier_4_attempts INT NOT NULL DEFAULT 0,
    xp INT NOT NULL DEFAULT 0
);

-- DROP TABLE BossData;
-- DROP TABLE SlayerBossType;
-- DROP TABLE SlayerEventData;
-- DROP TABLE GuildEvent;
-- DROP TABLE GuildEventType;
INSERT INTO GuildEventType (id, type_name)
VALUES
(0, 'slayer'),
(1, 'diana'),
(2, 'dungeons');
