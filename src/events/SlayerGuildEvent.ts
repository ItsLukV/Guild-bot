import { EmbedBuilder, Snowflake, TextChannel } from 'discord.js';
import { Bot } from '../Bot';
import { prettieEndOfSlayerEventEmbed } from '../utils/utilsFunction';
import { GuildEvent, GuildEventType } from './GuildEvent';
import { GuildUser } from '../utils/GuildUser';
import { Slayer } from '../@types/skyblockProfile';
import { todo } from '../utils/todo';

export type BossType = 'zombie' | 'spider' | 'wolf' | 'enderman' | 'blaze' | 'vampire';


export class SlayerGuildEvent extends GuildEvent {
    // public readonly eventType: GuildEventType;
    constructor(duration: number) {
        super(duration)
    }

    public getType(): GuildEventType {
        return "event_slayer";
    }

    public toString(): string {
        return `SlayerGuildEvent (id: ${this.getUUID()})`;
    }



    public async activate() {
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

    public async loadPlayer(player: GuildUser) {
        todo()
    }
}
