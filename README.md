# League of Legends Match Tracker

## Technologies

- **NestJS** - REST API & Riot API polling
- **Go** - Discord notifications + Local ARAM tracking
- **RabbitMQ** - Message queue
- **PostgreSQL** - Database
- **Docker** - Containerization
- **AWS** - Deployment (EC2)

## Architecture

```
                    CLOUD SERVICES
┌─────────────────────────────────────────────┐
│         NestJS Application (ECS)            │
│                                             │
│  • REST API for summoner management         │
│  • Riot API polling (ranked matches)        │
│  • ARAM game tracking (from LCU client)     │
│  • Hourly ranking aggregation               │
│                                             │
│         ┌───────────────────┐               │
│         │   PostgreSQL      │               │
│         └───────────────────┘               │
└─────────────────┬───────────────────────────┘
                  │
                  ▼
         ┌────────────────┐
         │   RabbitMQ     │
         │                │
         │ • summoners    │
         │ • matchs       │
         └────────┬───────┘
                  │
                  ▼
         ┌──────────────────┐
         │  Go Consumer     │
         │                  │
         │ Discord Notifier │
         └─────────┬────────┘
                   │
                   ▼
            Discord Server


                LOCAL (WINDOWS)
┌──────────────────────────────────────────────┐
│          Go LCU Client                       │
│                                              │
│  • Monitors League Client (friend list)      │
│  • Detects ARAM games                        │
│  • Reports to NestJS API (authenticated)     │
│                                              │
│         League Client (LCU)                  │
└──────────────────────────────────────────────┘
```

## How It Works

### Ranked Match Tracking (Cloud)
1. User adds summoner via NestJS API
2. NestJS polls Riot API every 30 seconds to check game status
3. When ranked match ends, NestJS publishes message to RabbitMQ
4. Go consumer receives message and sends Discord notification with match details (KDA, LP change, win/loss)

### ARAM Tracking (Local + Cloud)
1. Go LCU client runs on Windows, connects to local League Client
2. Authenticates with NestJS API (JWT)
3. Monitors friend list every 10 seconds for ARAM games
4. Reports game start/end to NestJS API
5. NestJS stores in PostgreSQL