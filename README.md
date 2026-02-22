# League of Legends Match Tracker

## Technologies

- **NestJS** - REST API & Riot API polling
- **Go** - Discord notification worker
- **RabbitMQ** - Message queue
- **PostgreSQL** - Database
- **Docker** - Containerization
- **AWS** - Deployment (ECS, RDS, Amazon MQ)

## Architecture
```
┌─────────────────────────────────────────────┐
│            NestJS Application               │
│                                             │
│  ┌─────────────┐      ┌─────────────────┐   │
│  │   REST API  │      │  Cron Scheduler │   │
│  │             │      │  (every 30s)    │   │
│  │ POST /track │      │                 │   │
│  │ DELETE /... │      │  Check Riot API │   │
│  └─────────────┘      └─────────┬───────┘   │
│                                 │           │
│         ┌───────────────────────▼─────────┐ │
│         │     PostgreSQL Database         │ │
│         │  (summoners, matches, state)    │ │
│         └───────────────────────┬─────────┘ │
│                                 │           │
│              Match detected? ───┘           │
│                     │                       │
└─────────────────────┼───────────────────────┘
                      │
                      ▼
            ┌─────────────────┐
            │   RabbitMQ      │
            │ match.notify    │
            └────────┬────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │  Go Consumer          │
         │  (Discord Notifier)   │
         │                       │
         │  - Consume messages   │
         │  - Format embed       │
         │  - Send to Discord    │
         │  - Handle retries     │
         └───────────┬───────────┘
                     │
                     ▼
              Discord Server
```

## How It Works

1. **User** submits summoner name via REST API
2. **NestJS** stores summoner in PostgreSQL
3. **Cron job** polls Riot API every 30s to check game status
4. **When match ends**, NestJS publishes message to RabbitMQ
5. **Go worker** consumes message and sends Discord notification