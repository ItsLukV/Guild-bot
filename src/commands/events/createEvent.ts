import { CacheType, ChatInputCommandInteraction, EmbedBuilder, Guild, SlashCommandBuilder } from "discord.js";
import { Command } from "../../@types/commands";
import { Bot } from "../../Bot";
import { GuildEvent, GuildEventType } from "../../utils/GuildEvent";

module.exports = {
    data: new SlashCommandBuilder()
    .setName('createevent')
    .setDescription('Creates a guild event')
    .addStringOption(option =>
        option.setName('event-type')
            .setDescription('The type of guild event')
            .setRequired(true)
            .addChoices(
        { name: 'slayer', value: 'event_slayer' },
                    { name: 'meow', value: 'event_meow' },
            ))
    .addIntegerOption(option =>
        option.setName('duration')
            .setDescription('Event duration in hours')
            .setRequired(true)
    )
            ,
    async execute(interaction: ChatInputCommandInteraction<CacheType>) {
        let event = new GuildEvent(
            interaction.options.getString('event-type', true) as GuildEventType,
            interaction.options.getInteger("duration",true)
        )
        Bot.guildEvents.push(event)
        let eventStatus = event.active(interaction)
        let embed = new EmbedBuilder()
        .setColor(eventStatus ? 0x00ff00 : 0xff0000 )
        .setDescription(eventStatus ? "Event created" : "Failed to create event")
        interaction.reply({ embeds: [embed] });
    }
} as Command;
