package mixer

// The mixer is an http server responsible for mixing the client coins by doing the following
// 1. On startup, preseed a certain amount of addresses with coins (to bootstrap the mixing process)
// 2. Accept requests from clients that conform to certain rules (min, max, fee, etc)
// 3. Provide a deposit address back to the client
// 4. Check the blockchain to see if/when the client sends funds to the deposit address
// 5. When funds are received, move funds into smaller random amounts into addresses mixer controls
// 6. Send those funds back to the clients specified address in random intervals
