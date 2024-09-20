import { Snowflake, VoiceState } from "discord.js";
import { UUID } from "./hypixel-api";

export class GuildUser {
    readonly id: Snowflake;
    readonly name: string;
    readonly uuid: string;

    constructor(id: Snowflake, name: string, uuid: UUID) {
        this.id = id;
        this.name = name
        this.uuid = uuid
    }
}
