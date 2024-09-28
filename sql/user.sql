CREATE TABLE users (
    id SERIAL INT PRIMARY KEY,
    discord_username text,
    discord_snowfale text,
    minecraft_username text,
    minecraft_uuid text
);

CREATE TABLE guild_event_type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE guild_event (
    id TEXT PRIMARY KEY,
    event_type INT REFERENCES guild_event_type(id),
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL
);

CREATE TABLE guild_event_entries (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    guild_id TEXT REFERENCES guild_event(id)
);

CREATE TABLE boss_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE boss_stats (
    id SERIAL PRIMARY KEY,
    boss_type_id INT REFERENCES boss_types(id),
    xp INT NOT NULL,
    tier_0_kills INT NOT NULL DEFAULT 0,
    tier_1_kills INT NOT NULL DEFAULT 0,
    tier_2_kills INT NOT NULL DEFAULT 0,
    tier_3_kills INT NOT NULL DEFAULT 0,
    tier_4_kills INT NOT NULL DEFAULT 0,
    tier_0_attempts INT NOT NULL DEFAULT 0,
    tier_1_attempts INT NOT NULL DEFAULT 0,
    tier_2_attempts INT NOT NULL DEFAULT 0,
    tier_3_attempts INT NOT NULL DEFAULT 0,
    tier_4_attempts INT NOT NULL DEFAULT 0
);

CREATE TABLE guild_event_boss_stats (
    id SERIAL PRIMARY KEY,
    guild_event_id TEXT REFERENCES guild_event(id),
    boss_stats_id INT REFERENCES boss_stats(id),
    user_id INT REFERENCES users(id)
);
