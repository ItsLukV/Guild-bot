import { CacheType, ChatInputCommandInteraction, EmbedBuilder, Guild, SlashCommandBuilder, TextChannel } from "discord.js";
import { Command } from "../../@types/commands";
import { Bot } from "../../Bot";
import { GuildEvent, GuildEventType } from "../../utils/GuildEvent";
import { GuildUser } from "../../utils/GuildUser";

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
                    // { name: 'meow', value: 'event_meow' },
            ))
    .addIntegerOption(option =>
        option.setName('duration')
            .setDescription('Event duration in hours')
            .setRequired(true)
    )
            ,
    async execute(interaction: ChatInputCommandInteraction<CacheType>) {
        if(interaction.memberPermissions?.has("Administrator") === false) {
            let embed = new EmbedBuilder()
            .setColor(0xff0000 )
            .setDescription("[ERROR] You dont have the right permissions")
            interaction.reply({ embeds: [embed] });
            interaction.reply({ embeds: [embed] });
            return
        }
        let event = new GuildEvent(
            interaction.options.getString('event-type', true) as GuildEventType,
            interaction.options.getInteger("duration",true)
        )
        let guildUser = Bot.getGuildMember(interaction.user.id)
        if (!guildUser) {
            let embed = new EmbedBuilder()
            .setColor(0xff0000 )
            .setDescription("[ERROR] Please register!")
            interaction.reply({ embeds: [embed] });
            return
        }
        let embed = await Bot.guildEventManager.createEvent(event);

        interaction.reply({ embeds: [embed] })

        await Bot.guildEventManager.addUser(event.getUUID(), new GuildUser({
            id: '300381646929133568',
            discordname: 'emilzacho',
            mcUsername: 'TTVEmilMZEU',
            uuid: '5ef04c7a95ae4c9396cefe925e4d5833'
        }));
        await Bot.guildEventManager.addUser(event.getUUID(), new GuildUser({
            id: '251350379650875394',
            discordname: 'rabbsdk',
            mcUsername: '22um',
            uuid: 'af1da1dcf5b046b3b412cc3af47f0bd6'
        }));
        await Bot.guildEventManager.addUser(event.getUUID(), new GuildUser({
            id: '428908206748729345',
            discordname: 'stepjepp',
            mcUsername: 'stepjeppe',
            uuid: '4c19fd1577234187b2d8d0fcea61e31e'
        }));

        await Bot.guildEventManager.addUser(event.getUUID(), guildUser);

        // TODO: Make this to a command
        await event.active()
    }
} as Command;
