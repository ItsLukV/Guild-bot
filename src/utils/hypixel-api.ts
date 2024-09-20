import { Profile, Skyblock } from "../@types/skyblockProfile";

type checkUser  = {player: {socialMedia?: {links?: {DISCORD?: string}}}}

export async function checkUser(uuid: UUID): Promise<string | null> {
    let path = `https://api.hypixel.net/v2/player?key=${process.env.HYPIXEL_API}&uuid=${uuid}`
    let content = await api<checkUser>(path)
    if (!content.player.socialMedia) {
        return null
    }
    if (!content.player.socialMedia.links) {
        return null
    }
    if (!content.player.socialMedia.links.DISCORD) {
        return null
    }
    return content.player.socialMedia.links.DISCORD
}

export type UUID = string

type minecraftApi = {
    id: UUID
    name: string
}



export async function getSkyblock(uuid: UUID): Promise<Skyblock> {
    let path = `https://api.hypixel.net/v2/skyblock/profiles?key=${process.env.HYPIXEL_API}&uuid=${uuid}`
    return await api<Skyblock>(path)
}

export async function getActiveProfile(uuid:UUID): Promise<Profile | undefined> {
    let skyblockData = await getSkyblock(uuid)
    return skyblockData.profiles?.find(profile => profile.selected)
}

export async function checkMCUsername(username:string): Promise<string | null>{
    return checkUser((await getMCUUID(username)).id)
}

export async function getMCUUID(username:string): Promise<minecraftApi> {
    return await api<minecraftApi>(`https://api.mojang.com/users/profiles/minecraft/${username}`)
}


async function api<T>(path: string): Promise<T> {
    const response = await fetch(path);
    if (!response.ok) {
        if (response.status === 403) {
            console.error("Forbidden: Check your API key or rate limits.");
        } else {
            console.error(`Error: ${response.statusText} (Status Code: ${response.status})`);
        }
        throw new Error(response.statusText);
    }
    return await response.json();
}
