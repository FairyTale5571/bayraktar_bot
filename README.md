# bayraktar_bot

## This discord bot developed for Arma 3 Rocket Life server.

## Build and depolyment

API: https://bot.rocket-rp.fun/swag/index.html#/

### Build and start

#### Requirements:
```
- Golang v1.18 or higher
- MariaDB v10.6 or higher
- MongoDB
- Redis
```

#### Environment variables
```bash
STEAMAPI_KEY=steam_api_key
DISCORD_TOKEN=discord_bot_token
VIP_ROLE_ID=VIP
GUILD_ID=123456789
REG_ROLE_ID=0
POST_PASSWORD=123456
MYSQL_URI=user:password@tcp(localhost:3306)/database
MONGO_URI=mongodb://localhost:27017
REDISCLOUD_URL=redis://localhost:6379
URL=http://localhost
PORT=3000
DEBUG=true
```

```bash
    make build
```
