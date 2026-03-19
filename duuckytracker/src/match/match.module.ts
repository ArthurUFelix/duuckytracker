import { Module } from '@nestjs/common';
import { MatchService } from './match.service';
import { MatchController } from './match.controller';
import { Match } from './entities/match.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { MatchSchedulerService } from './match-scheduler.service';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'MATCH_SERVICE',
        transport: Transport.RMQ,
        options: {
          urls: [String(process.env.RABBITMQ_URL)],
          queue: 'matches',
          queueOptions: {
            durable: true,
          },
        },
      },
    ]),
    TypeOrmModule.forFeature([Match]),
  ],
  controllers: [MatchController],
  providers: [MatchService, MatchSchedulerService],
})
export class MatchModule {}
