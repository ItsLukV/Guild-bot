use std::{env, vec};
use std::str::FromStr;
use serenity::async_trait;
use serenity::model::channel::Message;
use serenity::prelude::*;
use dotenv::dotenv;
use reqwest;
use serde_json::Value;
use serde::{Deserialize, Serialize};
use chrono::{DateTime, FixedOffset};

const prefix: &str = "!";
struct Handler;



#[derive(Debug, Deserialize, Serialize)]
pub struct GuildEvent {
    pub id: String,
    pub users: Vec<String>,
    #[serde(rename = "start_time")]
    pub start_time: DateTime<FixedOffset>,
    pub duration: u32,
    // #[serde(rename = "type")]
    // pub event_type: EventType,
    #[serde(rename = "is_hidden")]
    pub is_hidden: bool,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct ListUser {
    pub id: String,
    #[serde(rename = "active_profile_UUID")]
    pub active_profile_uuid: String,
    pub discord_snowflake: String,
    pub fetch_data: bool,
}
#[derive(Debug, Deserialize, Serialize)]
pub struct UsersResponse<T> {
    pub users: Vec<T>,
}

enum GuildEventTypes {
    Diana,
    Dungeons,
    Mining
}

impl FromStr for GuildEventTypes {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "diana" => Ok(GuildEventTypes::Diana),
            "dungeons" => Ok(GuildEventTypes::Dungeons),
            "mining" => Ok(GuildEventTypes::Mining),
            _ => Err(format!("'{}' is not a valid event type", s)),
        }
    }
}

enum Commands {
    Users,
    User(GuildEventTypes),
    GuildEvents(GuildEventTypes),
    List(ListType)
}

#[derive(Clone, Copy)]
enum ListType {
    Users,
    GuildEvents,
}

impl FromStr for ListType {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "users" => Ok(ListType::Users),
            "event" => Ok(ListType::GuildEvents),
            _ => Err(format!("{} is not a valid ListType", s))
        }
    }
}

impl ToString for ListType {
    fn to_string(&self) -> String {
        match self {
            ListType::Users => "users".to_owned(),
            ListType::GuildEvents => "events".to_owned(),
        }
    }
}

impl ListType {
    async fn get_list(self) -> Result<Vec<String>, String> {
        match self {
            ListType::Users => {
                let response = reqwest::get("https://lukv.dev/api/guildevents")
                    .await
                    .map_err(|e| e.to_string())?;

                println!("{:?}", response);


                let events: Vec<GuildEvent> = response.json()
                    .await
                    .map_err(|e| e.to_string())?;

                Ok(events.into_iter().map(|e| e.id).collect())
            },
            ListType::GuildEvents => {
                let response = reqwest::get("https://lukv.dev/api/users")
                .await
                .map_err(|e| e.to_string())?;


            let users_response: UsersResponse<ListUser> = response.json()
                .await
                .map_err(|e| e.to_string())?;

            Ok(users_response.users.into_iter().map(|u| u.active_profile_uuid).collect())
            }
        }
    }
}

impl FromStr for Commands {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        // Remove leading/trailing whitespace and split into words
        let command_str = if s.starts_with(prefix) {
            &s[prefix.len()..]
        } else {
            return Err(format!("Commands must start with '{}'", prefix));
        };

        let words: Vec<&str> = command_str.trim().split_whitespace().collect();

        if words.is_empty() {
            return Err("Empty command".to_string());
        }

        match words[0].to_lowercase().as_str() {
            "users" => {
                if words.len() > 1 {
                    Err("'users' command doesn't take arguments".to_string())
                } else {
                    Ok(Commands::Users)
                }
            },
            "user" => {
                if words.len() > 1 {
                    // Assuming you want to pass some user identifier
                    match GuildEventTypes::from_str(words[1]) {
                        Ok(v) => Ok(Commands::User(v)),
                        Err(v) => Err(v),
                    }
                } else {
                    Err("'user' command requires an argument".to_string())
                }
            },
            "event" => {
                if words.len() < 2 {
                    return Err("'event' command requires a type".to_string());
                }

                match GuildEventTypes::from_str(words[1]) {
                    Ok(v) => Ok(Commands::GuildEvents(v)),
                    Err(v) => Err(v),
                }
            },
            "list" => {
                if words.len() < 2 {
                    return Err("'event' command requires a type".to_string());
                }

                match ListType::from_str(words[1]) {
                    Ok(v) => Ok(Commands::List(v)),
                    Err(v) => Err(v),
                }
            }
            _ => Err(format!("'{}' is not a valid command", words[0])),
        }
    }
}

impl Commands {
    async fn run(&self, ctx: Context, msg: Message) {
        match self {
            Commands::Users => todo!(),
            Commands::User(guild_event_types) => todo!(),
            Commands::GuildEvents(guild_event_types) => todo!(),
            Commands::List(list_type) => {
                match list_type.get_list().await {
                    Ok(v) => {
                        let mut list = list_type.to_string();
                        for e in v {
                            if !list.is_empty() {
                                list.push('\n');
                                list.push('`');
                            }
                            list.push_str(&e);
                            list.push('`');

                        }
                        if let Err(e) = msg.channel_id.say(&ctx.http, list).await {
                            println!("Error sending message: {}", e);
                        }
                    },
                    Err(e) => println!("{}", e),
                }

            },
        }
    }
}

#[async_trait]
impl EventHandler for Handler {
    async fn message(&self, ctx: Context, message: Message) {
        match Commands::from_str(&message.content) {
            Ok(v) => v.run(ctx, message).await,
            Err(e) => println!("Failed to run command {}", e),
        }

    }
}

#[tokio::main]
async fn main() {
    dotenv().ok();

    let token = env::var("DISCORD_TOKEN").expect("Expected a token in the environment");
    let intents = GatewayIntents::GUILD_MESSAGES
        | GatewayIntents::DIRECT_MESSAGES
        | GatewayIntents::MESSAGE_CONTENT;

    let mut client = Client::builder(&token, intents)
        .event_handler(Handler)
        .await
        .expect("Err creating client");

    if let Err(why) = client.start().await {
        println!("Client error: {why:?}");
    }
}
