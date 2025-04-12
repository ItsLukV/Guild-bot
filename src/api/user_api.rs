use reqwest::{Error, StatusCode};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::vec;

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct RootUser {
    #[serde(rename = "diana_data")]
    pub diana_data: DianaData,
    #[serde(rename = "dungeons_data")]
    pub dungeons_data: DungeonsData,
    #[serde(rename = "mining_data")]
    pub mining_data: MiningData,
    pub user: User,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DianaData {
    pub id: String,
    #[serde(rename = "fetch_time")]
    pub fetch_time: String,
    #[serde(rename = "burrows_treasure")]
    pub burrows_treasure: i64,
    #[serde(rename = "burrows_combat")]
    pub burrows_combat: i64,
    #[serde(rename = "gaia_construct")]
    pub gaia_construct: i64,
    #[serde(rename = "minos_champion")]
    pub minos_champion: i64,
    #[serde(rename = "minos_hunter")]
    pub minos_hunter: i64,
    #[serde(rename = "minos_inquisitor")]
    pub minos_inquisitor: i64,
    pub minotaur: i64,
    #[serde(rename = "siamese_lynx")]
    pub siamese_lynx: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct DungeonsData {
    pub id: String,
    #[serde(rename = "fetch_time")]
    pub fetch_time: String,
    pub experience: f64,
    pub completions: Completions,
    #[serde(rename = "master_completions")]
    pub master_completions: MasterCompletions,
    #[serde(rename = "class_xp")]
    pub class_xp: ClassXp,
    pub secrets: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Completions {
    #[serde(rename = "0")]
    pub n0: i64,
    #[serde(rename = "1")]
    pub n1: i64,
    #[serde(rename = "2")]
    pub n2: i64,
    #[serde(rename = "3")]
    pub n3: i64,
    #[serde(rename = "4")]
    pub n4: i64,
    #[serde(rename = "5")]
    pub n5: i64,
    #[serde(rename = "6")]
    pub n6: i64,
    #[serde(rename = "7")]
    pub n7: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MasterCompletions {
    #[serde(rename = "1")]
    pub n1: i64,
    #[serde(rename = "2")]
    pub n2: i64,
    #[serde(rename = "3")]
    pub n3: i64,
    #[serde(rename = "4")]
    pub n4: i64,
    #[serde(rename = "5")]
    pub n5: i64,
    #[serde(rename = "6")]
    pub n6: i64,
    #[serde(rename = "7")]
    pub n7: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct ClassXp {
    pub archer: f64,
    pub berserk: f64,
    pub healer: f64,
    pub mage: f64,
    pub tank: f64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct MiningData {
    pub id: String,
    #[serde(rename = "fetch_time")]
    pub fetch_time: String,
    #[serde(rename = "mineshaft_count")]
    pub mineshaft_count: i64,
    #[serde(rename = "fossil_dust")]
    pub fossil_dust: i64,
    #[serde(rename = "tungsten_corpse")]
    pub tungsten_corpse: i64,
    #[serde(rename = "umber_corpse")]
    pub umber_corpse: i64,
    #[serde(rename = "lapis_corpse")]
    pub lapis_corpse: i64,
    #[serde(rename = "vanguard_corpse")]
    pub vanguard_corpse: i64,
    #[serde(rename = "nucleus_runs")]
    pub nucleus_runs: i64,
    #[serde(rename = "mithril_powder")]
    pub mithril_powder: i64,
    #[serde(rename = "powder_gemstone")]
    pub powder_gemstone: i64,
    #[serde(rename = "glacite_powder")]
    pub glacite_powder: i64,
    #[serde(rename = "Collections")]
    pub collections: Collections,
    #[serde(rename = "scatha_kills")]
    pub scatha_kills: i64,
    #[serde(rename = "worm_kills")]
    pub worm_kills: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Collections {
    pub mithril: i64,
    pub gemstone: i64,
    #[serde(rename = "gold_Ingot")]
    pub gold_ingot: i64,
    pub netherrack: i64,
    pub diamond: i64,
    pub ice: i64,
    #[serde(rename = "Redstone")]
    pub redstone: i64,
    #[serde(rename = "Lapis")]
    pub lapis: i64,
    pub sulphur: i64,
    pub coal: i64,
    pub emerald: i64,
    #[serde(rename = "end_stone")]
    pub end_stone: i64,
    pub glowstone: i64,
    pub gravel: i64,
    #[serde(rename = "iron_ingot")]
    pub iron_ingot: i64,
    pub mycelium: i64,
    pub quartz: i64,
    #[serde(rename = "Obsidian")]
    pub obsidian: i64,
    #[serde(rename = "red_sand")]
    pub red_sand: i64,
    pub sand: i64,
    pub cobblestone: i64,
    #[serde(rename = "hard_stone")]
    pub hard_stone: i64,
    #[serde(rename = "metal_heart")]
    pub metal_heart: i64,
    pub glacite: i64,
    pub umber: i64,
    pub tungsten: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct User {
    pub id: String,
    #[serde(rename = "active_profile_UUID")]
    pub active_profile_uuid: String,
    #[serde(rename = "discord_snowflake")]
    pub discord_snowflake: String,
    #[serde(rename = "fetch_data")]
    pub fetch_data: bool,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
struct ApiResponse {
    users: Vec<User>
}

// API functions
pub async fn fetch_users_data() -> Result<Vec<User>, String> {
    let response = reqwest::get("https://lukv.dev/api/users")
        .await
        .map_err(|e| format!("Failed to make request: {}", e))?;

    if !response.status().is_success() {
        return Err(format!("API request failed with status: {}", response.status()));
    }

    let api_response = response
        .json::<ApiResponse>()
        .await
        .map_err(|e| format!("Failed to parse response: {}", e))?;

    Ok(api_response.users)
}

pub async fn fetch_user_data(uuid: String) -> Result<RootUser, String> {
    let url = format!("https://lukv.dev/api/user?id={}", uuid);
    println!("Fetching user data from: {}", url);

    let response = reqwest::get(&url)
        .await
        .map_err(|e| format!("Request failed: {}", e))?;

    if !response.status().is_success() {
        let status = response.status();
        let error_msg = match status {
            StatusCode::NOT_FOUND => "User not found".to_string(),
            StatusCode::UNAUTHORIZED => "Authentication required".to_string(),
            _ => format!("API request failed with status: {}", status),
        };
        return Err(error_msg);
    }

    response.json::<RootUser>()
        .await
        .map_err(|e| format!("Failed to parse response: {}", e))
}

pub async fn uuid_to_username(uuid: &str) -> Option<String> {
    let clean_uuid = uuid.replace("-", "");
    let response = reqwest::get(&format!(
        "https://api.mojang.com/user/profile/{}",
        clean_uuid
    ))
    .await.ok()?
    .json::<serde_json::Value>()
    .await.ok()?;

    response.get("name")?.as_str().map(|s| s.to_string())
}

pub async fn username_to_uuid(username: &str) -> Option<String> {
    let url = format!(
        "https://api.mojang.com/users/profiles/minecraft/{}",
        username
    );

    reqwest::get(&url)
        .await.ok()?
        .json::<serde_json::Value>()
        .await.ok()?
        .get("id")?
        .as_str()
        .map(|s| s.to_string())
}

#[derive(Debug, Serialize, Deserialize)]
struct Guild {
    #[serde(rename = "_id")]
    id: String,
    // You can add other fields here if needed, but they'll be ignored
}

#[derive(Debug, Serialize, Deserialize)]
struct ResponseGuild {
    success: bool,
    guild: Guild,
}

pub async fn is_in_guild(uuid: String) -> Option<bool> {
    let api_token = std::env::var("HYPIXEL_API_KEY").ok()?;
    let api_url = format!(
        "https://api.hypixel.net/v2/guild?key={}&player={}",
        api_token, uuid
    );

    let response = reqwest::get(&api_url)
        .await.ok()?
        .json::<ResponseGuild>()
        .await.ok()?;

    Some(response.guild.id == "66b7cea08ea8c94d7358c510")
}

#[derive(Debug, Serialize, Deserialize)]
struct ResponseDiscord {
    success: bool,
    player: PlayerInfo,
}

#[derive(Debug, Serialize, Deserialize)]
struct PlayerInfo {
    #[serde(rename = "socialMedia")]
    social_media: SocialMedia
}

#[derive(Debug, Serialize, Deserialize)]
struct SocialMedia {
    #[serde(rename = "links")]
    links: SocialMediaLinks
}

#[derive(Debug, Serialize, Deserialize)]
struct SocialMediaLinks {
    #[serde(rename = "DISCORD")]
    discord: String
}

pub async fn get_discord(uuid: String) -> Option<String> {
    let api_token = std::env::var("HYPIXEL_API_KEY").ok()?;
    let api_url = format!(
        "https://api.hypixel.net/v2/player?key={}&uuid={}",
        api_token, uuid
    );

    let response = reqwest::get(&api_url)
        .await.ok()?
        .json::<ResponseDiscord>()
        .await.ok()?;

    Some(response.player.social_media.links.discord)
}
