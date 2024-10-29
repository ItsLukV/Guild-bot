import { Embed, EmbedBuilder } from "discord.js";
import { Slayer } from "../@types/skyblockProfile";
import { Bot } from "../Bot";
import { GuildUser } from "../utils/GuildUser";
import { getActiveProfile, getSkyblock } from "../utils/hypixel-api";
import { GuildEvent } from "./GuildEvent";
import { todo } from "../utils/todo";



export class GuildEventManager {

    private readonly guildEvents: Map<string, GuildEvent> = new Map();

    async createEvent(event: GuildEvent): Promise<EmbedBuilder> {
        this.guildEvents.set(event.getUUID(),event);
        const client = await Bot.pool.connect()

        const eventId = await client.query(`
            INSERT INTO guild_event (id, event_type)
            VALUES ($1, (SELECT id FROM guild_event_type WHERE name = $2))
            RETURNING id;
        `, [event.getUUID(), event.getType()]);
        return new EmbedBuilder()
        .setTitle( "Successfully created a " + event)
        // .setDescription(`Ends in <t:${Math.round(new Date().getTime() / 1000) + event.duration}:R>`)
        .addFields()
        .setColor(0x00ff00)
    }

    async addUser(eventId: string, user: GuildUser): Promise<EmbedBuilder> {
        let embed;
        try {
            let event = this.guildEvents.get(eventId)
            if(!event) {
                throw new Error("Missing event")
            }
            let playerData = await getActiveProfile(user.uuid)
            let slayerData = playerData?.members[user.uuid]?.slayer;
            if(slayerData) {
                embed = new EmbedBuilder().setTitle("success!")
            } else {
                throw new Error("No slayer data")
            }
            event.addUser(user);
        } catch (e) {
            console.error(e)
            embed = new EmbedBuilder().setTitle("error please contact please a admin")
        }
        return embed;
    }


    async load() {
        await this.loadEvent();
        // await this.loadPlayerData();
    }

    async loadPlayerData() {
        todo()
    }

    async activeEvent(eventId: string): Promise<EmbedBuilder> {
        let event = this.guildEvents.get(eventId)
        if(!event) {
            return new EmbedBuilder()
                .setTitle('Failed to start the event')
                .setColor(0xff0000)
        }

        const client = await Bot.pool.connect()
        try {
            const _ = await client.query(`
                UPDATE guild_events
                SET start_date = NOW(), end_date = NOW() + INTERVAL '${event.duration} seconds
                WHERE id = ${eventId};
            `);
        } finally {
            client.release()
        }

        // Return immediately since setTimeout is non-blocking
        return new EmbedBuilder()
        .setTitle( "Successfully created a" + event)
        .setDescription(`Ends in <t:${Math.round(new Date().getTime() / 1000) + event.duration}:R>`)
        .addFields()
        .setColor(0x00ff00)
    }

    async loadEvent() {
        let query = `
        SELECT
            ge.id AS event_id,
            ge.start_date,
            ge.end_date,
            ge.end_date - ge.start_date as duration,
            get.name AS event_type,
            u.id AS user_id,
            u.discord_username,
            u.minecraft_username,
            u.minecraft_uuid
        FROM
            guild_event ge
        JOIN
            guild_event_type get ON ge.event_type = get.id
        LEFT JOIN
            guild_event_entries gee ON ge.id = gee.guild_id
        LEFT JOIN
            users u ON gee.user_id = u.id;
        `

        const client = await Bot.pool.connect()
        try {
            const res = await client.query(query)
            console.log(res)
            // for (let row of res.rows) {
            //     if (!this.guildEvents.has(row.event_id)) {
            //         this.guildEvents.set(
            //             row.event_id,
            //             new GuildEvent(
            //                 row.duration
            //             )
            //         )
            //     }

            //     this.guildEvents.get(row.event_id)!
            //         .addUser(new GuildUser({
            //             id: row.user_id,
            //             discordname: row.discord_username,
            //             uuid: row.minecraft_uuid,
            //             mcUsername: row.minecraft_username
            //         }))
            // }
        } finally {
            client.release()
        }
    }
}
