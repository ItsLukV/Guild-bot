import {v4 as uuidv4} from 'uuid';
import { Bot } from '../Bot';
import { ChatInputCommandInteraction, CacheType, EmbedBuilder } from 'discord.js';
import { checkMCUsername, getActiveProfile, getMCUUID, getSkyblock } from './hypixel-api';

export type GuildEventType = "event_slayer" | "event_meow"

export class GuildEvent {
    private readonly _uuid;
    private user: GuildEvent[] = [];
    private guildType: GuildEventType;
    private duration: number;
    private finshed = false;

    constructor(guildType: GuildEventType, duration: number) {
        console.log(`Created a guildEventType: ${guildType} with a duration of ${duration}`)
        this.guildType = guildType;
        this.duration = duration
        this._uuid = uuidv4()
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
            try {
                let user = Bot.guildMembers.get(interaction.user.id)
                if (!user) {
                    throw new Error("User not registed")

                }
                let uuid = user?.id
                if (!uuid) {
                    console.log("Failed to retrieve UUID");
                }
                let profileData  = await getActiveProfile(uuid);
                if (profileData) {
                    console.log(profileData.members[uuid]?.slayer);
                } else {
                    console.log("Failed to retrieve Profile");
                }
            } catch (error) {
                console.log("An error occurred:", error);
            }
        }, this.duration * 1_000); // this.duration should be in seconds

        // Return immediately since setTimeout is non-blocking
        return false;
    }




}
