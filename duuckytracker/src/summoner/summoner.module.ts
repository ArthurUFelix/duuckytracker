import { Module } from '@nestjs/common';
import { SummonerService } from './summoner.service';
import { SummonerController } from './summoner.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Summoner } from './entities/summoner.entity';
import { LeagueApiService } from './league-api/league-api.service';
import { SummonerSchedulerService } from './summoner-scheduler.service';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'SUMMONER_SERVICE',
        transport: Transport.RMQ,
        options: {
          urls: [
            `amqp://${process.env.RABBITMQ_HOST}:${process.env.RABBITMQ_PORT}`,
          ],
          queue: 'summoners',
          queueOptions: {
            durable: true,
          },
        },
      },
    ]),
    TypeOrmModule.forFeature([Summoner]),
  ],
  controllers: [SummonerController],
  providers: [SummonerService, LeagueApiService, SummonerSchedulerService],
})
export class SummonerModule {}
