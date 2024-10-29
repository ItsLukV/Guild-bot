import { CacheType, ChatInputCommandInteraction, SlashCommandBuilder } from "discord.js";
import { Command } from "../../@types/commands";

module.exports = {
    data: new SlashCommandBuilder()
    .setName('activateevent')
    .setDescription('Starts/Activates a guild event')
    .addStringOption(option =>
        option.setName('event-type')
            .setDescription('The type of guild event')
            .setRequired(true)
            .addChoices(
        { name: 'slayer', value: 'event_slayer' },
                    // { name: 'meow', value: 'event_meow' },
            ))
    ,
    async execute(interaction: ChatInputCommandInteraction<CacheType>) {

    }
} as Command;
