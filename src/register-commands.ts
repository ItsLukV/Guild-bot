import { REST, Routes } from 'discord.js';
import fs from 'node:fs';
import path from 'node:path';
import dotenv from 'dotenv';

dotenv.config();

const commands: any[] = [];
// Grab all the items (files/folders) from the commands directory
const commandsPath: string = path.join(__dirname, 'commands');
const commandItems: string[] = fs.readdirSync(commandsPath);

for (const item of commandItems) {
	// Create the full path to the item
	const itemPath: string = path.join(commandsPath, item);

	// Check if the item is a directory or a file
	const isDirectory: boolean = fs.lstatSync(itemPath).isDirectory();

	if (isDirectory) {
		// If it's a folder, grab the command files inside it
		const commandFiles: string[] = fs.readdirSync(itemPath).filter((file: string) => file.endsWith('.js'));
		for (const file of commandFiles) {
			const filePath: string = path.join(itemPath, file);
			const command = require(filePath);
			if ('data' in command && 'execute' in command) {
				commands.push(command.data.toJSON());
			} else {
				console.log(`[WARNING] The command at ${filePath} is missing a required "data" or "execute" property.`);
			}
		}
	} else if (item.endsWith('.js')) {
		// If it's a file, add the command directly
		const command = require(itemPath);
		if ('data' in command && 'execute' in command) {
			commands.push(command.data.toJSON());
		} else {
			console.log(`[WARNING] The command at ${itemPath} is missing a required "data" or "execute" property.`);
		}
	}
}

// Construct and prepare an instance of the REST module
const rest = new REST().setToken(process.env.DISCORD_TOKEN!);

// and deploy your commands!
(async () => {
	try {
		console.log(`Started refreshing ${commands.length} application (/) commands.`);

		// The put method is used to fully refresh all commands in the guild with the current set
		const data: any = await rest.put(
			Routes.applicationGuildCommands(process.env.CLIENT_ID!, process.env.GUILD_ID!),
			{ body: commands },
		);

		console.log(`Successfully reloaded ${data.length} application (/) commands.`);
	} catch (error) {
		// And of course, make sure you catch and log any errors!
		console.error(error);
	}
})();
