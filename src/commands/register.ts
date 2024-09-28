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
        const userName = interaction.options.getString("ign", true);
        const embed = new EmbedBuilder()

        // Check if Minecraft username is linked to a Discord username
        const hypixelDiscordUsername = await checkMCUsername(userName);

        if (!hypixelDiscordUsername) {
            embed.setDescription("Please add your Discord username on Hypixel.");
            embed.setColor(0xff0000);
        } else if (Bot.hasGuildMember(interaction.user.id)) {
            embed.setDescription("You are already registered!");
            embed.setColor(0xff0000);
        } else if (hypixelDiscordUsername === interaction.user.username) {
            const mcUUID = await getMCUUID(userName);

            Bot.addGuildMember(
                interaction.user.id,
                new GuildUser(
                    interaction.user.id,
                    interaction.user.username,
                    mcUUID.id,
                    userName
                )
            );

            console.log(`Registered: ${interaction.user.username}`);
            embed.setColor(0x00ff00);
            embed.setDescription("You are now registered!");
        } else {
            embed.setColor(0xff0000);
            embed.setDescription("Your Discord username does not match the one on Hypixel.");
        }

        await interaction.reply({ embeds: [embed] });
    }

} as Command;
