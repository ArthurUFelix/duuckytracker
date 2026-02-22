import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { SummonerModule } from './summoner/summoner.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Summoner } from './summoner/entities/summoner.entity';
import { ScheduleModule } from '@nestjs/schedule';
import { UserModule } from './user/user.module';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      url: process.env.DATABASE_URL,
      entities: [Summoner],
      synchronize: true,
    }),
    SummonerModule,
    ScheduleModule.forRoot(),
    ConfigModule.forRoot(),
    UserModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule {}
