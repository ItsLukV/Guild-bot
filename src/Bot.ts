import { Guild, IntentsBitField, Snowflake } from "discord.js";
import { GuildUser } from "./utils/GuildUser";
import { Discord } from "./Discord";
import { GuildEvent } from "./utils/GuildEvent";

export class Bot {

  public static guildEvents: GuildEvent[] = [];

  public static readonly guildMembers: Map<Snowflake, GuildUser> = new Map();
  public readonly discord = new Discord({
    allowedMentions: { parse: ['users', 'roles'], repliedUser: true },
    intents: [
        IntentsBitField.Flags.Guilds,
        IntentsBitField.Flags.GuildMessages,
    ],
  });
    static Discord: any;

  public async start() {
    await this.discord.start()
  }

}
