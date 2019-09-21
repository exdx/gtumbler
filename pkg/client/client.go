package client

// The client is an http client that sends requests to the mixer consisting of the following
// 1. The initial request sends a list of new addresses that the mixed coins will eventually be sent back to
// 2. The mixer responds with a deposit address
// 3. The client sends the full deposit amount to the deposit address
// 4. From that point on the client checks the list of addresses sent in (1) to be notified when their mixing coins are available
