# Ares

A Moderator for Slack.

## Installation

- Create a new Slack App with following permissions: `admin`, `bot`, `channels:history`, `channels:read`, `channels:write`, `chat:write:bot`, `files:read`, `files:write:user`, `groups:read`, `groups:write`, `users:read`.
- Install the bot on [Heroku](https://www.heroku.com/deploy/?template=https://github.com/avinassh/ares)

## Features

- Adds the bot to all channels at initialization 
- Deletes uploaded images and reuploads them to Imgur
- Mutes an user (type in `ares mute <username>`)
- Moderators

### Moderators

Ares can make some members as moderators. To add moderators, make sure you have set `MOD_IDS` env variable with the comma separated user ids of the mods, like `U1AQSSBSA,U0CM1JMV5,U0SDDB26B`

Currently, moderators can remove or mute other users.

## License

The mighty MIT license. Please check `LICENSE` for more details.
