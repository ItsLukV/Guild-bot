mod api;

use std::thread::sleep;
use std::time::Duration;

use api::user_api::{fetch_user_data, get_discord, is_in_guild, username_to_uuid};
use api::user_api::{fetch_users_data, uuid_to_username};
use poise::serenity_prelude::{json, CreateEmbedFooter};
use poise::{serenity_prelude as serenity, ChoiceParameter};
use poise::CreateReply;
use reqwest::StatusCode;
use separator::Separatable;
use serde_json::json;
use once_cell::sync::Lazy;
use ::serenity::futures::io::Repeat;

static DEBUG: Lazy<String> = Lazy::new(|| {
    // Can compute this at runtime
    std::env::var("DEBUG").unwrap_or("0".to_string())
});

static API_URL: Lazy<String> = Lazy::new(|| {
    // Can compute this at runtime
    std::env::var("API_URL").expect("Failed to get API_URL")
});

static API_TOKEN: Lazy<String> = Lazy::new(|| {
    // Can compute this at runtime
    std::env::var("API_TOKEN").expect("Failed to get API_TOKEN")
});

struct Data {} // User data
type Error = Box<dyn std::error::Error + Send + Sync>;
type Context<'a> = poise::Context<'a, Data, Error>;

/// Get list of users with their Minecraft names
#[poise::command(slash_command, prefix_command)]
async fn get_users(
    ctx: Context<'_>,
) -> Result<(), Error> {
    ctx.defer().await?;

    match fetch_users_data().await {
        Ok(users) if !users.is_empty() => {
            let mut message = String::new();



            for (i, user) in users.iter().enumerate() {
                let limit = 3;
                if i > limit {
                    message.push_str(&format!("- {} More Users", users.len() - limit));
                    break;
                }


                if i > 0 {
                    sleep(Duration::from_millis(50));
                }

                let username = uuid_to_username(&user.id)
                    .await
                    .unwrap_or_else(|| "Unknown".to_string());

                message.push_str(&format!("- {}: {}\n", user.id, username));
            }


            let embed = serenity::CreateEmbed::default()
                .title("Users")
                .description(message)
                .color(serenity::Colour::DARK_GREEN);

            // Correct way to send a reply with content and reply setting
            ctx.send(CreateReply::default()
                .embed(embed)
                .reply(true)
            ).await?;
        }
        Err(_) => {
            let embed = serenity::CreateEmbed::default()
                .title("No Users Found")
                .color(serenity::Colour::RED);

            ctx.send(CreateReply::default()
                .embed(embed)
            ).await?;
        }
        _ => {
            let embed = serenity::CreateEmbed::default()
                .title("API Error")
                .description("Failed to fetch user data")
                .color(serenity::Colour::RED);

            ctx.send(CreateReply::default()
                .embed(embed)
            ).await?;
        }
    }

    Ok(())
}

fn format_decimal(n: f64) -> String {
    format!("{:.2}", n)
        .parse::<f64>()
        .unwrap()
        .separated_string()
}

#[derive(Debug, ChoiceParameter)]
enum UserDataType{
    #[name = "mining"]
    Mining,
    #[name = "diana"]
    Diana,
    #[name = "dungeons"]
    Dungeons,
    #[name = "status"]
    Status
}

#[poise::command(slash_command, prefix_command)]
async fn get_user(
    ctx: Context<'_>,
    #[description = "The Minecraft username to look up"] username: String,
    #[description = "The Minecraft username to look up"] data: UserDataType,
) -> Result<(), Error> {
    ctx.defer().await?;
    println!("{:?}", username_to_uuid(&username.clone().as_str()).await);
    let Some(uuid) = username_to_uuid(&username.clone().as_str()).await else {
        let embed = serenity::CreateEmbed::default()
            .title("API Error")
            .description("Not a vaild Minecraft username")
            .color(serenity::Colour::RED);

        ctx.send(CreateReply::default()
            .embed(embed)
        ).await?;
        return Ok(());
    };

    if let Ok(user) = fetch_user_data(uuid).await {
        let footer = CreateEmbedFooter::new("Use /link to connect your account");

        let embed = serenity::CreateEmbed::default()
                .title(format!("{}'s Profile", username))
                .color(serenity::Colour::DARK_GREEN)
                .thumbnail(format!("https://mc-heads.net/avatar/{}/100.png", user.user.id))
                .fields(
                    match data {
                        UserDataType::Mining => {vec![
                            ("Mithril Powder",      format!("`{}`", user.mining_data.mithril_powder.separated_string()),       true),
                            ("Gemstone Powder",     format!("`{}`", user.mining_data.powder_gemstone.separated_string()),      true),
                            ("Glacite Powder",      format!("`{}`", user.mining_data.glacite_powder.separated_string()),       false),
                            ("Lapis Corpse",        format!("`{}`", user.mining_data.lapis_corpse.separated_string()),         true),
                            ("Tungsten Corpse",     format!("`{}`", user.mining_data.tungsten_corpse.separated_string()),      true),
                            ("Umber Corpse",        format!("`{}`", user.mining_data.umber_corpse.separated_string()),         true),
                            ("Vanguard Corpse",     format!("`{}`", user.mining_data.vanguard_corpse.separated_string()),      false),
                            ("Gemstone Collection", format!("`{}`", user.mining_data.collections.gemstone.separated_string()), true),
                            ("Mithril Collection",  format!("`{}`", user.mining_data.collections.mithril.separated_string()),  true),
                            ("Glacite Collection",  format!("`{}`", user.mining_data.collections.glacite.separated_string()),  true),
                            ("Tungsten Collection", format!("`{}`", user.mining_data.collections.tungsten.separated_string()), true),
                            ("Umber Collection",    format!("`{}`", user.mining_data.collections.umber.separated_string()),    true),
                        ]},
                        UserDataType::Diana => {vec![
                            ("Treasure Burrows",    format!("`{}`", user.diana_data.burrows_treasure.separated_string()),  false),
                            ("Combat Burrows",      format!("`{}`", user.diana_data.burrows_combat.separated_string()),    false),
                            ("Minos Inquisitor",    format!("`{}`", user.diana_data.minos_inquisitor.separated_string()),  false),
                            ("Minos Champion",      format!("`{}`", user.diana_data.minos_champion.separated_string()),    false),
                            ("Minotaur",            format!("`{}`", user.diana_data.minotaur.separated_string()),          false),
                            ("Gaia Construct",      format!("`{}`", user.diana_data.gaia_construct.separated_string()),    false),
                            ("Siamese_lynx",        format!("`{}`", user.diana_data.siamese_lynx.separated_string()),      false),
                            ("Minos Hunter",        format!("`{}`", user.diana_data.minos_hunter.separated_string()),      false),

                        ]},
                        UserDataType::Dungeons => {vec![
                            ("Experience",  format!("`{}`", format_decimal(user.dungeons_data.experience)),              false),
                            ("Secrets",     format!("`{}`", user.dungeons_data.secrets.separated_string()),                 false),
                            ("F1 Comp",     format!("`{}`", user.dungeons_data.completions.n1.separated_string()),          true),
                            ("F2 Comp",     format!("`{}`", user.dungeons_data.completions.n2.separated_string()),          true),
                            ("F3 Comp",     format!("`{}`", user.dungeons_data.completions.n3.separated_string()),          true),
                            ("F4 Comp",     format!("`{}`", user.dungeons_data.completions.n4.separated_string()),          true),
                            ("F5 Comp",     format!("`{}`", user.dungeons_data.completions.n5.separated_string()),          true),
                            ("F6 Comp",     format!("`{}`", user.dungeons_data.completions.n6.separated_string()),          true),
                            ("F7 Comp",     format!("`{}`", user.dungeons_data.completions.n7.separated_string()),          false),
                            ("M1 Comp",     format!("`{}`", user.dungeons_data.master_completions.n1.separated_string()),   true),
                            ("M2 Comp",     format!("`{}`", user.dungeons_data.master_completions.n2.separated_string()),   true),
                            ("M3 Comp",     format!("`{}`", user.dungeons_data.master_completions.n3.separated_string()),   true),
                            ("M4 Comp",     format!("`{}`", user.dungeons_data.master_completions.n4.separated_string()),   true),
                            ("M5 Comp",     format!("`{}`", user.dungeons_data.master_completions.n5.separated_string()),   true),
                            ("M6 Comp",     format!("`{}`", user.dungeons_data.master_completions.n6.separated_string()),   true),
                            ("M7 Comp",     format!("`{}`", user.dungeons_data.master_completions.n7.separated_string()),   false),
                            ("Archer Xp",   format!("`{}`", format_decimal(user.dungeons_data.class_xp.archer)),            true),
                            ("Berserk Xp",  format!("`{}`", format_decimal(user.dungeons_data.class_xp.berserk)),           true),
                            ("Healer Xp",   format!("`{}`", format_decimal(user.dungeons_data.class_xp.healer)),            true),
                            ("Mage Xp",     format!("`{}`", format_decimal(user.dungeons_data.class_xp.mage)),              true),
                            ("Tank Xp",     format!("`{}`", format_decimal(user.dungeons_data.class_xp.tank)),              true),
                        ]},
                        _ => {vec![
                                ("Last Fetch for Diana Data", user.diana_data.fetch_time.clone(), true),
                                ("Last Fetch for Dungeons Data", user.dungeons_data.fetch_time.clone(), true),
                                ("Last Fetch for Dungeons Data", user.dungeons_data.fetch_time.clone(), true),
                                ("Account Status", if user.user.fetch_data { "✅ Active".to_string() } else { "❌ Inactive".to_string() }, true),
                            ]}
                    }
                ).footer(footer);


        ctx.send(CreateReply::default()
            .embed(embed)
        ).await?;
    } else {
        let embed = serenity::CreateEmbed::default()
            .title("API Error")
            .description("Make sure the username is correct")
            .color(serenity::Colour::RED);

        ctx.send(CreateReply::default()
            .embed(embed)
        ).await?;
    }

    Ok(())
}

#[poise::command(slash_command, prefix_command)]
async fn link(
    ctx: Context<'_>,
    #[description = "The Minecraft username to look up"] username: String,
) -> Result<(), Error> {
    let Some(uuid) = username_to_uuid(&username.clone().to_string()).await else {
        let embed = serenity::CreateEmbed::default()
        .title("API Error")
        .description("Failed to get UUID")
        .color(serenity::Colour::RED);

        ctx.send(CreateReply::default()
        .embed(embed)).await?;

        return Ok(());
    };


    match is_in_guild(uuid.clone()).await {
        None => {
            let embed = serenity::CreateEmbed::default()
            .title("API Error")
            .description(format!("The user `{}` is not part of Skyblock Locals", username))
            .color(serenity::Colour::RED);

            ctx.send(CreateReply::default()
            .embed(embed)).await?;

            return Ok(());
        }
        _ => println!("User is in the guild")
    }

    if let Some(discord_name) = get_discord(uuid.clone()).await {
        let e = ctx.author().name.as_ref();
        if e != discord_name.to_lowercase() {
            let embed = serenity::CreateEmbed::default()
            .title("API Error")
            .description("Discord names are not the same")
            .color(serenity::Colour::RED);

            ctx.send(CreateReply::default()
            .embed(embed)).await?;
        } else {
            // Create a client
            let client = reqwest::Client::new();
            let url = &format!("{}/user?token={}",API_URL.to_string() ,API_TOKEN.to_string());
            println!("{}", url.clone());
            // Send a POST request
            let response = client
                .post(url)
                .json(&json!({
                    "uuid": uuid,
                    "discord_snowflake": ctx.author().id
                }))
                .send()
                .await?;

                match response.status() {
                    StatusCode::CREATED => {
                        let embed = serenity::CreateEmbed::default()
                        .title("Linked")
                        .description(":)")
                        .color(serenity::Colour::DARK_GREEN);

                        ctx.send(CreateReply::default()
                        .embed(embed)).await?;
                    },
                    StatusCode::INTERNAL_SERVER_ERROR => {
                        println!("{}", format!("INTERNAL_SERVER_ERROR for url: {}", url));

                        let embed = serenity::CreateEmbed::default()
                        .title("API Error")
                        .description("Something went wrong")
                        .color(serenity::Colour::RED);

                        ctx.send(CreateReply::default()
                        .embed(embed)).await?;
                    },
                    _ => {
                        println!("{}", format!("INTERNAL_SERVER_ERROR for url: {}", url));

                        let embed = serenity::CreateEmbed::default()
                        .title("API Error")
                        .description("Something went wrong")
                        .color(serenity::Colour::RED);

                        ctx.send(CreateReply::default()
                        .embed(embed)).await?;
                    }
                }


        }
    } else  {
        let embed = serenity::CreateEmbed::default()
        .title("API Error")
        .description("Please link your discord name to your Hypixel account")
        .color(serenity::Colour::RED);

        ctx.send(CreateReply::default()
        .embed(embed)).await?;
    }

    Ok(())
}

#[tokio::main]
async fn main() {
    let token = std::env::var("DISCORD_TOKEN").expect("missing DISCORD_TOKEN");
    let intents = serenity::GatewayIntents::non_privileged();

    let framework = poise::Framework::builder()
        .options(poise::FrameworkOptions {
            commands: vec![get_users(), get_user(), link()],
            ..Default::default()
        })
        .setup(|ctx, _ready, framework| {
            Box::pin(async move {
                // Check DEBUG environment variable
                let debug_mode = std::env::var("DEBUG")
                    .map(|s| s == "true" || s == "1")
                    .unwrap_or(false);

                if debug_mode {
                    // Register only in test guild for debugging
                    let test_guild_id = serenity::GuildId::new(1355176342915645600); // Your test server ID
                    poise::builtins::register_in_guild(
                        ctx,
                        &framework.options().commands,
                        test_guild_id
                    ).await?;
                    println!("Registered commands in test guild only (DEBUG mode)");
                } else {
                    // Register globally in production
                    poise::builtins::register_globally(
                        ctx,
                        &framework.options().commands
                    ).await?;
                    println!("Registered commands globally");
                }

                Ok(Data {})
            })
        })
        .build();

    let mut client = serenity::ClientBuilder::new(token, intents)
        .framework(framework)
        .await
        .expect("Error creating client");

    client.start().await.expect("Error running client");
}
