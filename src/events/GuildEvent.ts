import { Snowflake } from "discord.js";
import { GuildUser } from "../utils/GuildUser";
import { v4 as uuidv4 } from 'uuid';


export type GuildEventType = "event_slayer" | "event_meow"

export abstract class GuildEvent {
    private readonly _uuid;
    public readonly duration: number;
    protected finshed = false;
    protected users: Map<Snowflake, GuildUser> = new Map()


    constructor(duration: number) {
        console.log(`Created a guildEventType: ${typeof(this)} with a duration of ${duration} seconds.`)
        this.duration = duration;
        this._uuid = uuidv4();
    }

    public abstract active(): void;

    public getUser(id: string): GuildUser | undefined{
        return this.users.get(id)
    }

    public addUser(user: GuildUser) {
        this.users.set(user.id, user)
    }

    public getUUID(): string {
        return this._uuid;
    }

    public abstract getType(): GuildEventType;
    public abstract toString(): string;
}
