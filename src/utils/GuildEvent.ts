import {v4 as uuidv4} from 'uuid'
import { Bot } from '../Bot';
import { ChatInputCommandInteraction, CacheType, EmbedBuilder, Snowflake, TextChannel } from 'discord.js';
import { checkMCUsername, getActiveProfile, getMCUUID, getSkyblock } from './hypixel-api';
import { GuildUser } from './GuildUser';
import { PlayerData, Slayer } from '../@types/skyblockProfile';
import { todo } from './todo';
import { prettieEndOfSlayerEventEmbed } from './utilsFunction';
import { BossType } from './GuildEventManager';

export type GuildEventType = "event_slayer" | "event_meow"

export class GuildEvent {
    private readonly _uuid;
    private user: GuildEvent[] = [];
    public readonly eventType: GuildEventType;
    public readonly duration: number;
    private finshed = false;
    private users: Map<Snowflake, GuildUser> = new Map()



    constructor(guildType: GuildEventType, duration: number) {
        console.log(`Created a guildEventType: ${guildType} with a duration of ${duration}`)
        this.eventType = guildType;
        this.duration = duration
        this._uuid = uuidv4()
    }

    public getUser(id: string): GuildUser | undefined{
        return this.users.get(id)
    }

    public addUser(user: GuildUser) {
        this.users.set(user.id, user)
    }

    public getUUID(): string {
        return this._uuid;
    }

    public async active() {
        let embed = new EmbedBuilder().setTitle("Events started")
        .setDescription(`Ends in <t:${Math.round(new Date().getTime() / 1000) + this.duration}:R>`)

        let channel = Bot.discord.channels.cache.get(process.env.GUILD_DISCORD_CHANNEL as string);
        // Check if the channel is a TextChannel
        if (channel instanceof TextChannel) {
            await channel.send({ embeds: [embed] });
        } else {
            console.error("The channel is not a text channel or does not exist.");
        }

        setTimeout(async () => {
            let embed = await prettieEndOfSlayerEventEmbed([...this.users.values()])
            let channel = Bot.discord.channels.cache.get(process.env.GUILD_DISCORD_CHANNEL as string);
            // Check if the channel is a TextChannel
            if (channel instanceof TextChannel) {
                await channel.send({ embeds: [embed] });
            } else {
                console.error("The channel is not a text channel or does not exist.");
            }
        }, this.duration * 1_000);
    }
}

export function prettieEventType(type: GuildEventType): string {
    switch (type) {
        case "event_slayer":
            return "Slayer Event"
        default:
            return type
    }
}
