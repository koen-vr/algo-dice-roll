# algo-roller
 
Demo's an on-chain dice roller with simple pseudo random numbers.
 
This demo stores private keys in plain text and should not be used in production.
 
This project **hasn't been security audited** and should not be used in a production environment.
 
## Brief
 
Algo roller implements a simple pseudo random number generator where the outcomes are predictable and verifiable based on the seed used to initialize its state. The demo simulates a simple dice rolling game between two players.
 
## Requirements
 
* Linux or macOS
* Golang version 1.17.0 or higher
* Python 3. The scripts assumes the Python executable is called `python3`.
* The [Algorand Node software](https://developer.algorand.org/docs/run-a-node/setup/install/). A private network is used, hence there is no need to sync up MainNet or TestNet. `goal` is assumed to be in the PATH.
 
## Setup
 
To install all required packages, run:
```bash
python3 -m pip install -r requirements.txt
```
 
## Usage
 
TODO: Implement
 
=====
 
Why roll pseudo random numbers on-chain? To verify the outcome of an encounter once the results are being published and before rewards are paid out.
 
In practice when an event or encounter starts:
 
1) Players hash their initial state and publish that hash.
2) The published hash is stored on-chain and keeps the state private.
3) ... the event or encounter plays out ...
4) At the end the initial state, action logs and results are published.
5) With the state published the app can now verify it by checking the hash.
6) Once done, the state can be used to replay the actions and verify the result.
