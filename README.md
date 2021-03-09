# POWWOW

"Word of Wisdom" tcp server with the Prof of Work DDOS protection.

## Startup

You can run server and 100 instances of test clients with the
following command:

```sh
make compose
```

## Details

* A custom implementation of hashcash was chosen as the PoW algorithm.
* As hashing function for hashcash sha3-256 was selected.
* Nonce complexity is selected on the fly by the current server CPU load.
