# Ares

A Moderator for Slack.

## Installation

- Create a new Slack App with following permissions: `admin`, `bot`, `channels:history`, `channels:read`, `channels:write`, `chat:write:bot`, `files:read`, `files:write:user`, `groups:read`, `groups:write`, `users:read`.
- Install the bot on [Heroku](https://www.heroku.com/deploy/?template=https://github.com/avinassh/ares)

## Features

- Adds the bot to all channels at initialization 
- Deletes uploaded images and reuploads them to Imgur
- Mutes an user (type in `ares mute <username>`)
- [Moderators](mods.md)

## License

The mighty MIT license. Please check `LICENSE` for more details.
