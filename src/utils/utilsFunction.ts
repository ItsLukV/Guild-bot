import { EmbedBuilder, Snowflake } from "discord.js";
import { Slayer, SlayerBoss } from "../@types/skyblockProfile";
import { GuildUser } from "./GuildUser";
import { getActiveProfile } from "./hypixel-api";

export async function prettieEndOfSlayerEventEmbed(users: GuildUser[]): Promise<EmbedBuilder> {

    let data = await calcSlayerEvent(users)
    let txt = ""
    let index = 1;
    data.forEach((value,key) =>{
        txt += `${index}: <@${key}>: ${value} points\n`
        index++;
    });

    return new EmbedBuilder()
    .setTitle('Slayer Event Results') // Optional, but good to include
    .setDescription('Here are the results from the slayer event.') // Ensure this is populated
    .addFields({
        name: "Leaderboard",
        value: txt
    })
    .setColor('#0099ff');
}

export async function calcSlayerEvent(users: GuildUser[]) {
    let usersMap = new Map<Snowflake, number>();
    for (let user of users) {
        let oldSlayerDatas = user.playerStats.slayer.slayer_bosses;
        let data = 0
        let apiReq = await getActiveProfile(user.uuid)
        let newSlayerData = apiReq!.members[user.uuid]!.slayer.slayer_bosses

        for (let slayerData of Object.keys(newSlayerData) as (keyof typeof newSlayerData)[]) {
            const oldSlayerBossData = oldSlayerDatas[slayerData];
            const newSlayerBossData = newSlayerData[slayerData];
            let pointsInfo = slayerPoints[slayerData]
            data += ((newSlayerBossData.boss_attempts_tier_0 || 0) - (oldSlayerBossData.boss_attempts_tier_0 || 0) ) * pointsInfo.tier0
            data += ((newSlayerBossData.boss_attempts_tier_1 || 0) - (oldSlayerBossData.boss_attempts_tier_1 || 0) ) * pointsInfo.tier1
            data += ((newSlayerBossData.boss_attempts_tier_2 || 0) - (oldSlayerBossData.boss_attempts_tier_2 || 0) ) * pointsInfo.tier2
            data += ((newSlayerBossData.boss_attempts_tier_3 || 0) - (oldSlayerBossData.boss_attempts_tier_3 || 0) ) * pointsInfo.tier3
            data += ((newSlayerBossData.boss_attempts_tier_4 || 0) - (oldSlayerBossData.boss_attempts_tier_4 || 0) ) * pointsInfo.tier4
        }
        usersMap.set(user.id,data)
    }

    return new Map([...usersMap.entries()].sort((a, b) => b[1] - a[1]));
}


const slayerPoints = {
    zombie: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
    spider: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
    wolf: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
    enderman: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
    blaze: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
    vampire: {
        tier0: 1,
        tier1: 2,
        tier2: 3,
        tier3: 4,
        tier4: 5,
        xp: 1
    },
}
