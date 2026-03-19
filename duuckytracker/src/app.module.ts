import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SummonerModule } from './summoner/summoner.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Summoner } from './summoner/entities/summoner.entity';
import { ScheduleModule } from '@nestjs/schedule';
import { AuthModule } from './auth/auth.module';
import { User } from './user/entities/user.entity';
import { UserModule } from './user/user.module';
import { MatchModule } from './match/match.module';
import { Match } from './match/entities/match.entity';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      url: process.env.DATABASE_URL,
      entities: [Summoner, User, Match],
      synchronize: true,
    }),
    ScheduleModule.forRoot(),
    ConfigModule.forRoot(),
    SummonerModule,
    AuthModule,
    UserModule,
    MatchModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
