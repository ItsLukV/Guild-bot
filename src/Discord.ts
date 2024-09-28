import { Client, ClientOptions, Collection, Events } from 'discord.js';
import fs from 'node:fs';
import path from 'node:path';
import { Command } from './@types/commands';
import { Pool } from 'pg';

export class Discord extends Client {
    public commands: Collection<string, Command>;



    constructor(options: ClientOptions) {
        super(options)
        this.commands = new Collection();
    }

    private async registerCommands() {
        this.commands = new Collection();
        const foldersPath: string = path.join(__dirname, 'commands');
        const commandItems: string[] = fs.readdirSync(foldersPath);

        for (const item of commandItems) {
            const itemPath: string = path.join(foldersPath, item);
            const stat = fs.lstatSync(itemPath);

            // If the item is a directory, handle it as a folder of commands
            if (stat.isDirectory()) {
                const commandFiles: string[] = fs.readdirSync(itemPath).filter((file: string) => file.endsWith('.js'));
                for (const file of commandFiles) {
                    const filePath: string = path.join(itemPath, file);
                    const command = require(filePath);

                    if ('data' in command && 'execute' in command) {
                        this.commands.set(command.data.name, command);
                    } else {
                        console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
                    }
                }
            }
            // If the item is a file, handle it directly
            else if (stat.isFile() && item.endsWith('.js')) {
                const command = require(itemPath);

                if ('data' in command && 'execute' in command) {
                    this.commands.set(command.data.name, command);
                } else {
                    console.log(`[WARNING] The command at ${itemPath} is missing a required "data" or "execute" property.`);
                }
            }
        }
    }


    public async start() {
        await this.discordLogin();
        await this.registerCommands();

        this.on(Events.InteractionCreate, async (interaction) => {
            if (!interaction.isChatInputCommand()) return;

            const command = (interaction.client as Discord).commands.get(interaction.commandName);

            if (!command) {
                console.error(`No command matching ${interaction.commandName} was found.`);
                return;
            }

            try {
                await command.execute(interaction);
            } catch (error) {
                console.error(error);
                if (interaction.replied || interaction.deferred) {
                    await interaction.followUp({ content: 'There was an error while executing this command!', ephemeral: true });
                } else {
                    await interaction.reply({ content: 'There was an error while executing this command!', ephemeral: true });
                }
            }
        });
    }

    private async discordLogin() {
        await this.login(process.env.DISCORD_TOKEN);
        if (this.user) {
            console.log(`Logged in as ${this.user.tag}!`);
        } else {
            console.error("Failed to login");
        }
    }
}
