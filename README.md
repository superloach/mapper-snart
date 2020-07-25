![Docker](https://github.com/superloach/mapper/workflows/Docker/badge.svg)

# mapper
A Discord bot for Ingress, Pokémon GO, and Harry Potter: Wizards Unite.

## docker build
This repository is a plugin for [Snart](https://github.com/go-snart/snart). Thus, an easy way to build a Docker image for a bot with this plugin is using [go-snart/example](https://github.com/go-snart/example).
```sh
# in github.com/go-snart/example
./genplug github.com/superloach/mapper
docker build -t mapper .
```
