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

    public active(interaction: ChatInputCommandInteraction<CacheType>): boolean {
        if (!interaction.member) {
            interaction.reply({embeds: [new EmbedBuilder().setDescription("Error").setColor(0xff0000)]});
            return false;
        }
        setTimeout(async () => {
            let embed = await prettieEndOfSlayerEventEmbed([...this.users.values()])
            let channel = Bot.discord.channels.cache.get('1283422887482626153');
            // Check if the channel is a TextChannel
            if (channel instanceof TextChannel) {
                await channel.send({ embeds: [embed] });
            } else {
                console.error("The channel is not a text channel or does not exist.");
            }
        }, this.duration * 1_000);

        // Return immediately since setTimeout is non-blocking
        return true;
    }
}
