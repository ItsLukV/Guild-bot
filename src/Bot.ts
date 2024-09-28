import { GatewayIntentBits, Guild, IntentsBitField, Partials, Snowflake } from "discord.js";
import { GuildUser } from "./utils/GuildUser";
import { Discord } from "./Discord";
import { GuildEvent } from "./utils/GuildEvent";
require('dotenv').config()
import { Pool } from "pg";
import { todo } from "./utils/todo";
import { GuildEventManager } from "./utils/GuildEventManager";



export class Bot {
  public static pool: Pool = new Pool({
    user: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    host: process.env.DB_HOST,
    port: process.env.DB_PORT ? parseInt(process.env.DB_PORT) : undefined,
    database: process.env.DB,
  });

  public static readonly discord = new Discord({
    allowedMentions: { parse: ['users', 'roles'], repliedUser: true },
    intents: [
        IntentsBitField.Flags.Guilds,
        IntentsBitField.Flags.GuildMessages,
    ],
    partials: [
      Partials.Channel
    ]
  });



  public static readonly guildEventManager: GuildEventManager = new GuildEventManager()
  private static readonly guildMembers: Map<Snowflake, GuildUser> = new Map();

  public static async addGuildMember(key: Snowflake, value: GuildUser) {
    if (this.guildMembers.has(key)) {
      throw new Error("The user is already registed")
    }

    const client =  await Bot.pool.connect()
    try {
      const query = `INSERT INTO users (discord_username, discord_snowfale, minecraft_username, minecraft_uuid)
                       VALUES ($1, $2, $3, $4) RETURNING *`;
      const values = [value.discordName, value.id, value.mcUsername, value.uuid ]
        const res = client.query(query, values, (err: any, result: { rows: any[]; }) => {
          if (err) {
            console.error('Error executing query', err);
          }
        });
    } finally {
      client.release();
    }

    this.guildMembers.set(key, value)
  }

  public static hasGuildMember(key: Snowflake): boolean {
    return this.guildMembers.has(key)
  }

  public static getGuildMember(key: Snowflake): GuildUser | undefined {
    return this.guildMembers.get(key)
  }

  public async loadUsers() {
    const client = await Bot.pool.connect()

    try {
      const query = `SELECT * FROM users`;
      const res = await client.query(query);
      if (res.rows.length > 0) {
        for (let row of res.rows) {
          Bot.guildMembers.set(row.discord_snowfale, new GuildUser(
            row.discord_snowfale,
            row.discord_username,
            row.minecraft_uuid,
            row.minecraft_username
          ));
        }
      } else {
        console.log('No data found');
      }
    } finally {
      client.release();
    }
  }

  public async start() {
    await this.loadUsers()
    await Bot.guildEventManager.load()
    await Bot.discord.start()
  }
}
