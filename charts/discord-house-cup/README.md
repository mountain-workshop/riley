[![Build Status](https://drone.iratepublik.com/api/badges/sudermans/discord-house-cup/status.svg)](https://drone.iratepublik.com/sudermans/discord-house-cup)

# Discord House Cup

A discord bot to keep a leaderboard of points for a set of teams.

## Deployment

Drone is used to build and push a docker image when a tag is created. That tagged image is then used with the helm chart to deploy the newest version.
