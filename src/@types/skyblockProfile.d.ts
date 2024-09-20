export interface Skyblock {
    success: boolean,
    profiles: Profile[]
}

interface Profile {
    profile_id: string,
    community_upgrades: any,
    members: Record<string, PlayerData>
    selected: boolean
}

interface PlayerData {
    rift: Rift,
    player_data: PlayerData,
    glacite_player_data: any,
    events: any
    garden_player_data: any
    accessory_bag_storage: any
    leveling: any
    item_data: any
    jacobs_contest: any
    currencies: any
    dungeons: any
    profile: any
    pets_data: any
    player_id: string
    nether_island_player_data: any
    experimentation: any
    mining_core: any
    bestiary: any
    quests: any
    player_stats: any
    winter_player_data: any
    forge: any
    fairy_soul: any
    slayer: Slayer
    trophy_fish: any
    objectives: any
    inventory: any
    shared_inventory: any
    collection: any
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
    boss_kills_tier_0: number,
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
