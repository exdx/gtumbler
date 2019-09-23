# gtumbler

gtumbler is a JobCoin mixer written in Go. 

## Structure

gtumbler follows a client server architecture. 

### client
The client is an http client that sends requests to the mixer consisting of the following
1. The initial request sends a list of new addresses that the mixed coins will eventually be sent back to
2. The mixer responds with a deposit address
3. The client sends the full deposit amount to the deposit address
4. From that point on the client checks the list of addresses sent in (1) to be notified when their mixing coins are available

### mixer
The mixer is an http server responsible for mixing the client coins by doing the following
1. On startup, preseed a certain amount of addresses with coins (to bootstrap the mixing process).
Since there is no API method for creating coins from scratch, these addresses need to be made from the UI.
The mixer will assume these addresses already exist and are funded
2. Accept requests from clients that conform to certain rules (min, max, fee, etc)
3. Provide a deposit address back to the client
4. Check the blockchain to see if/when the client sends funds to the deposit address
5. When funds are received, move funds into smaller random amounts into addresses mixer controls
6. Send those funds back to the clients specified address

The mixer is the core focus of the project. Inside of the mixer is an additional service called the tumbler which helps 
the mixer fulfill its responsibilities. The tumbler is responsible for the actual mixing process. 

### tumbler
Given an array of addresses and a deposit address the tumbler does the following
1. Split the amount in the deposit address into random sizes
2. Allocate the random sizes into house array of addresses 
3. Check if others are also mixing and mix their coins in as well
3. Report the initial tumbling step as complete

for sending the coins back out to the customer the tumbler must do the opposite
1. Given an array of house addresses it must send the appropriate amount of coins to the array of customer accounts
2. It must send them in random sizes 
3. Report the final tumbling process as complete 

The way that the tumbler achieves randomness is through a predefined set of strategies that determine in what chunks the coins 
are going to be split. The strategies are an array containing different ways of chunking up an amount of cryptocurrency.
For example one strategy is [0.5, 0.5] which represents cutting up the deposit amount into halves. Ideally the strategies would not be predetermined but determined at runtime, 
but this comes at the cost of additional complexity, performance, and rounding errors. So for now only a set number of strategies are used by the tumbler.
The strategies are limited by the number of pre-seeded house accounts (5)

## Install and run

### Recommended: Docker

`bash docker-build.sh` installs both the client and mixer images on the local docker daemon. 

Run a server container in one terminal, and the client in another. _Start the server before the client_. 

_Terminal 1_: `docker run --rm -p 8989:8989 gtumbler/mixer:v0.0.1`

_Terminal 2_: `docker run --rm -p 8989:8989 gtumber/client:v0.0.1`

### Running locally
Be sure to have a local Go 1.11+ environment setup with support for go modules. Run `bash build.sh` to generate client and mixer binaries. 

Run the server in one terminal, and the client in another. _Start the server before the client_.

_Terminal 1_: `./gtumbler-mixer`

_Terminal 2_: `./gtumbler-client`

## Testing

Tests were written for the mixer only because of time constraints. Several unit tests were written. 

To run tests locally run `go test ./...`

## Positives

There is a good amount of resuse of code between the client and server. Nearly all packages import the crypto library which holds
semantics related to creating addresses and sending cryptocurrency. 

Custom types: even though the addresses and amounts are represented as strings in the API in the code these are custom types:
`crypto.Address` and `crypto.Amount` are used instead of string types to improve readability.

The interaction between the tumbler and the mixer is interesting. The mixer is responsible for the high-level handling off client requests
whereas the tumbler is concerned with splitting funds across addresses in a pseudo-random way. This is a nice separation of concerns. 

## Improvements

Initially the tumbler was envisioned to be more complex, by offsetting transactions over certain random times as well as chunks.
This would improve privacy and is not very hard to add but was cut because of time constraints. It also would generate more advanced stratgies
using randomness instead of a pre-determined number of strategies. 

The server does not inform the client in case of errors. For example, if the client sends in 200 coins, whereas the limit is 10, this information 
should be relayed back to the client and the coins should be sent back. Not difficult to add, but cut in the interest of time. 

The client would ideally be more of a true CLI instead of simply a tool that runs once and exits. Since the client is not really a core 
part of the project this was cut.  

 

