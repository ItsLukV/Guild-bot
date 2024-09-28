export interface Skyblock {
    success: boolean,
    profiles: Profile[]
}

interface Profile {
    profile_id: string,
    community_upgrades: unknown ,
    members: Record<string, PlayerData>
    selected: boolean
}

interface PlayerData {
    rift: Rift,
    glacite_player_data: unknown ,
    events: unknown
    garden_player_data: unknown
    accessory_bag_storage: unknown
    leveling: unknown
    item_data: unknown
    jacobs_contest: unknown
    currencies: unknown
    dungeons: unknown
    profile: unknown
    pets_data: unknown
    player_id: string
    nether_island_player_data: unknown
    experimentation: unknown
    mining_core: unknown
    bestiary: unknown
    quests: unknown
    player_stats: unknown
    winter_player_data: unknown
    forge: unknown
    fairy_soul: unknown
    slayer: Slayer
    trophy_fish: unknown
    objectives: unknown
    inventory: unknown
    shared_inventory: unknown
    collection: unknown
}

interface Slayer {
    slayer_quest: {
        type: string,
        tier: number,
        start_timestamp: number,
        completion_state: 0,
        used_amour: boolean,
        solo: boolean
    },
    slayer_bosses: {
        zombie: SlayerBoss,
        spider: SlayerBoss,
        wolf: SlayerBoss,
        enderman: SlayerBoss,
        blaze: SlayerBoss,
        vampire: SlayerBoss,
    }
}

interface SlayerBoss {
    claimed_levels: {
        level_1?: boolean,
        level_2?: boolean,
        level_3?: boolean,
        level_4?: boolean
        level_5?: boolean
        level_6?: boolean
        level_7_special?: boolean
        level_8_special?: boolean
        level_9_special?: boolean
    },
    boss_attempts_tier_0?: number,
    boss_kills_tier_0?: number,
    xp?: number,
    boss_attempts_tier_1?: number,
    boss_kills_tier_1?: number,
    boss_attempts_tier_2?: number,
    boss_kills_tier_2?: number,
    boss_attempts_tier_3?: number,
    boss_kills_tier_3?: number,
    boss_attempts_tier_4?: number,
    boss_kills_tier_4?: number
}
