import { Inject, Injectable, Logger } from '@nestjs/common';
import { Cron } from '@nestjs/schedule';
import { ClientProxy } from '@nestjs/microservices';
import { MatchService } from './match.service';
import { MoreThan } from 'typeorm';
import { randomUUID } from 'crypto';

@Injectable()
export class MatchSchedulerService {
  private readonly logger = new Logger(MatchSchedulerService.name);

  constructor(
    private readonly matchService: MatchService,
    @Inject('MATCH_SERVICE') private client: ClientProxy,
  ) {}

  @Cron('0 22 * * *')
  async checkMatchesRanking() {
    this.logger.log('Checking matches ranking...');

    try {
      const oneHourAgo = new Date();
      oneHourAgo.setHours(oneHourAgo.getHours() - 1);
      const matches = await this.matchService.findAll({
        where: {
          date: MoreThan(oneHourAgo),
        },
      });

      const matchesBySummoner = Object.groupBy(matches, (i) => i.summonerName);
      const matchesRanking = Object.entries(matchesBySummoner)
        .map(([summoner, matchesList]) => ({
          summoner,
          matchesCount: matchesList?.length ?? 0,
        }))
        .sort((a, b) => b.matchesCount - a.matchesCount);

      if (matchesRanking.length === 0) return;

      this.client.emit('send_matches', {
        id: randomUUID(),
        data: {
          matchesRanking,
        },
      });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      this.logger.error(`Error in checkMatchesRanking: ${message}`);
    }
  }
}
