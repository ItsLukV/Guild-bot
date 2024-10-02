import { Snowflake } from "discord.js";
import { UUID } from "./hypixel-api";
import { PlayerData, Slayer, SlayerBoss } from "../@types/skyblockProfile";

export class GuildUser {

    readonly id: Snowflake;
    readonly discordName: string;
    readonly uuid: string;
    readonly mcUsername: string;

    public playerStats: PlayerData = createPlayerDataFiller()

    constructor(
        idOrInput: Snowflake | { id: Snowflake, discordname: string, uuid: UUID, mcUsername: string },
        discordName?: string,
        uuid?: UUID,
        mcUsername?: string
    ) {
        if (typeof idOrInput === 'string') {
            // Case 1: Individual parameters
            this.id = idOrInput;
            this.discordName = discordName!;
            this.uuid = uuid!;
            this.mcUsername = mcUsername!;
        } else {
            // Case 2: Object input
            this.id = idOrInput.id;
            this.discordName = idOrInput.discordname;
            this.uuid = idOrInput.uuid;
            this.mcUsername = idOrInput.mcUsername;
        }
    }

    addPlayerStats(data: PlayerData) {
        this.playerStats = data;
    }

    addSlayer(data: Slayer) {
        this.playerStats.slayer = data;
    }

}

class PlayerStats {
    public slayer: Slayer;
    constructor(data: PlayerData) {
        this.slayer = data.slayer
    }
}


function createPlayerDataFiller(): PlayerData {
    return {
        slayer: {
            slayer_quest: {
                type: "",
                tier: 0,
                start_timestamp: 0,
                completion_state: 0,
                used_amour: false,
                solo: false
            },
            slayer_bosses: {
                zombie: { claimed_levels: {} },
                spider: { claimed_levels: {} },
                wolf: { claimed_levels: {} },
                enderman: { claimed_levels: {} },
                blaze: { claimed_levels: {} },
                vampire: { claimed_levels: {} }
            }
        },
        rift: undefined,
        glacite_player_data: undefined,
        events: undefined,
        garden_player_data: undefined,
        accessory_bag_storage: undefined,
        leveling: undefined,
        item_data: undefined,
        jacobs_contest: undefined,
        currencies: undefined,
        dungeons: undefined,
        profile: undefined,
        pets_data: undefined,
        player_id: "",
        nether_island_player_data: undefined,
        experimentation: undefined,
        mining_core: undefined,
        bestiary: undefined,
        quests: undefined,
        player_stats: undefined,
        winter_player_data: undefined,
        forge: undefined,
        fairy_soul: undefined,
        trophy_fish: undefined,
        objectives: undefined,
        inventory: undefined,
        shared_inventory: undefined,
        collection: undefined
    }
}
