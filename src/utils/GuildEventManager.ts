import { Slayer } from "../@types/skyblockProfile";
import { Bot } from "../Bot";
import { GuildEvent } from "./GuildEvent";
import { GuildUser } from "./GuildUser";
import { getActiveProfile, getSkyblock } from "./hypixel-api";

export type BossType = 'zombie' | 'spider' | 'wolf' | 'enderman' | 'blaze' | 'vampire';


export class GuildEventManager {

    private readonly guildEvents: Map<string, GuildEvent> = new Map();

    async createEvent(event: GuildEvent) {
        this.guildEvents.set(event.getUUID(),event);
        const client = await Bot.pool.connect()

        const eventId = await client.query(`
            INSERT INTO guild_event (id, event_type, start_date, end_date)
            VALUES ($1, (SELECT id FROM guild_event_type WHERE name = $2), NOW(), NOW() + INTERVAL '${event.duration} seconds')
            RETURNING id;
        `, [event.getUUID(), event.eventType]);
    }

    async addUser(eventId: string, user: GuildUser) {
        let event = this.guildEvents.get(eventId)
        if(!event) {
            throw new Error("Missing event")
        }

        switch(event.eventType) {
            case "event_slayer":
                let playerData = await getActiveProfile(user.uuid)
                let slayerData = playerData?.members[user.uuid]?.slayer;
                if(slayerData) {
                    user.addSlayer(slayerData);
                }
        }


        event.addUser(user);
    }


    async load() {
        await this.loadEvent();
        await this.loadPlayerData();
    }

    async loadPlayerData() {
        let query = `
            SELECT
                gebs.guild_event_id,
                u.id AS user_id,
                bt.name AS boss_type,
                bs.xp,
                bs.tier_0_kills,
                bs.tier_1_kills,
                bs.tier_2_kills,
                bs.tier_3_kills,
                bs.tier_4_kills,
                bs.tier_0_attempts,
                bs.tier_1_attempts,
                bs.tier_2_attempts,
                bs.tier_3_attempts,
                bs.tier_4_attempts
            FROM
                guild_event_boss_stats gebs
            JOIN
                users u ON gebs.user_id = u.id
            JOIN
                boss_stats bs ON gebs.boss_stats_id = bs.id
            JOIN
                boss_types bt ON bs.boss_type_id = bt.id
            ORDER BY guild_event_id, user_id;
        `
        const client = await Bot.pool.connect()
        try {
            const res = await client.query(query)
            let data: Slayer = {
                slayer_quest: {
                    type: "",
                    tier: 0,
                    start_timestamp: 0,
                    completion_state: 0,
                    used_amour: false,
                    solo: false
                },
                slayer_bosses: {
                    zombie: {claimed_levels: {}},
                    spider: {claimed_levels: {}},
                    wolf: {claimed_levels: {}},
                    enderman: {claimed_levels: {}},
                    blaze: {claimed_levels: {}},
                    vampire: {claimed_levels: {}}
                }
            };
            let user;
            for (let row of res.rows) {
                let event = this.guildEvents.get(row.guild_event_id);
                if (!event) {
                    throw new Error("Missing guild event")
                }

                if (!user) {
                    user = event.getUser(row.user_id)
                }

                if (!user) {
                    throw new Error("Missing user")
                }

                const bossTypeKey = row.boss_type as BossType;

                if (!(bossTypeKey in data.slayer_bosses)) {
                    throw new Error(`Unknown boss type: ${row.boss_type}`);
                }

                const bossData = data.slayer_bosses[bossTypeKey];

                bossData.boss_kills_tier_0 = row.tier_0_kills || 0;
                bossData.boss_kills_tier_1 = row.tier_1_kills || 0;
                bossData.boss_kills_tier_2 = row.tier_2_kills || 0;
                bossData.boss_kills_tier_3 = row.tier_3_kills || 0;
                bossData.boss_kills_tier_4 = row.tier_4_kills || 0;

                bossData.boss_attempts_tier_0 = row.tier_0_attempts || 0;
                bossData.boss_attempts_tier_1 = row.tier_1_attempts || 0;
                bossData.boss_attempts_tier_2 = row.tier_2_attempts || 0;
                bossData.boss_attempts_tier_3 = row.tier_3_attempts || 0;
                bossData.boss_attempts_tier_4 = row.tier_4_attempts || 0;

                bossData.xp = row.xp || 0;

                if (row.user_id === user.id) {
                    user.addSlayer(data);
                } else {
                    user = event.getUser(row.user_id)
                }
            }
        } finally {
            client.release()
        }
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
            for (let row of res.rows) {
                if (!this.guildEvents.has(row.event_id)) {
                    this.guildEvents.set(
                        row.event_id,
                        new GuildEvent(
                            row.event_type,
                            row.duration
                        )
                    )
                }

                this.guildEvents.get(row.event_id)!
                    .addUser(new GuildUser({
                        id: row.user_id,
                        discordname: row.discord_username,
                        uuid: row.minecraft_uuid,
                        mcUsername: row.minecraft_username
                    }))
            }
        } finally {
            client.release()
        }
    }
}
