# BF4DB-Search-Tool
Simple tool for searching for Players on the BF4DB database.


### Prerequisites
- A valid API key from BF4DB. You can get one [here](https://bf4db.com/patreon "here").


### Usage
To search for a player, run the following command:

```shell
$ bf4db <player name>
```
This will search for the player's information on BF4DB and print their information. If the player is not found on BF4DB, no information is printed.

### Options
- -h, --help: Show the help message.
- -c, --config: Set a new BF4DB API key.
- dbg: Show debug information.

### Example
```shell
$ bf4db PlayerName
Searching for PlayerName...

PlayerName | Not reported | Cheat score = X | BF4DB link: https://bf4db.com/player/ID/PlayerID
 Cheat Report: http://bf4cheatreport.com/PlayerID
```

# Contributions are welcome!
We welcome and appreciate contributions from everyone. Whether you're interested in fixing a bug, adding a new feature, or suggesting improvements.

To get started, make sure you have Go 1.16 or later installed on your system. Then, follow these steps:
1. Fork this repository to your own GitHub account.
2. Clone your fork of the repository to your local machine.
3. Make your changes and commit them with clear commit messages.
4. Push your changes to your forked repository on GitHub.
5. Submit a pull request to our repository.

If you're not comfortable submitting code changes, you can still help by opening an issue or feature request [here](https://github.com/pruu-airlines/BF4DB-Search-Tool/issues) :)
