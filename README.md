# mapper
Ingress, Pokémon Go, and Wizards Unite Discord bot.

## docker build
This repository is a plugin for [Snart](https://github.com/go-snart/snart). Thus, an easy way to build a Docker image for a bot with this plugin is using [go-snart/example](https://github.com/go-snart/example).
```sh
# in github.com/go-snart/example
./genplug github.com/superloach/mapper
docker build -t mapper .
```
