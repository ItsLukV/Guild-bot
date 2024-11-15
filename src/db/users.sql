CREATE TABLE IF NOT EXISTS users (
     discord_snowflake text PRIMARY KEY,
     discord_username text,
     minecraft_username text,
     minecraft_uuid text
);