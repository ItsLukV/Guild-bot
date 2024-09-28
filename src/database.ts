// const { Pool } = require('pg');

// const dbConfig = {
// 	user: 'postgres',
// 	password: '2002',
// 	host: 'localhost',
// 	port: '5432',
// 	database: 'guild',
// };


// const client = new Pool(dbConfig);

// // Connect to the database
// client
// 	.connect()
// 	.then(() => {
// 		console.log('Connected to PostgreSQL database');

// 		// Insert data into the Users table
// 		const query = `INSERT INTO "Users" (discord_username, discord_snowfale, minecraft_username, minecraft_uuid)
//                        VALUES ($1, $2, $3, $4) RETURNING *`;

// 		const values = ['user1', 'snowflake123', 'minecraftUser1', 'uuid1234']; // Replace with your actual data


// 		// Execute SQL queries here

// 		client.query(query, values, (err: any, result: { rows: any[]; }) => {
// 			if (err) {
// 				console.error('Error executing query', err);
// 			} else {
// 				console.log('Data inserted:', result.rows[0]);
// 			}

// 			// Close the connection when done
// 			client
// 				.end()
// 				.then(() => {
// 					console.log('Connection to PostgreSQL closed');
// 				})
// 				.catch((err: any) => {
// 					console.error('Error closing connection', err);
// 				});
// 		});
// 	})
// 	.catch((err: any) => {
// 		console.error('Error connecting to PostgreSQL database', err);
// 	});
