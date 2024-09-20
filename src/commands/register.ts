import { CacheType, ChatInputCommandInteraction, EmbedBuilder, SlashCommandBuilder } from "discord.js";
import { GuildUser } from "../utils/GuildUser";
import { checkMCUsername, getMCUUID } from "../utils/hypixel-api";
import { Discord } from "../Discord";
import { Command } from "../@types/commands";
import { Bot } from "../Bot";

module.exports = {
	data: new SlashCommandBuilder()
		.setName('register')
		.setDescription('Link your minecraft account with your discord account')
        .addStringOption(option =>
			option
            .setName('ign')
            .setRequired(true)
            .setDescription('Your minecraft username')),
    async execute(interaction: ChatInputCommandInteraction<CacheType>) {
        let userName = interaction.options.get("ign",true).value as string
        let embed = new EmbedBuilder().setColor(0x00ff00);
        let hypixelDiscordUsername = await checkMCUsername(userName);
        if(hypixelDiscordUsername === null) {
            embed.setDescription("Please add your discord username on hypixel");
        } else if (Bot.guildMembers.has(interaction.user.id)) {
            embed.setDescription("You are already registered!");
        } else if (hypixelDiscordUsername === interaction.user.username) {
            Bot.guildMembers.set(
                interaction.user.id,
                new GuildUser(
                    interaction.user.id,
                    interaction.user.username,
                    (await getMCUUID(userName)).id
                )
            );
            console.log("Registered: " + interaction.user.username);
            embed.setDescription("You are now registered!");
        } else {
            embed.setDescription("Tour discord username does not match the one on hypixel");
        }
        interaction.reply({ embeds: [embed] });
	},
} as Command;
