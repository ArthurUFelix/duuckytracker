import { Inject, Injectable, Logger } from '@nestjs/common';
import { Cron } from '@nestjs/schedule';
import { SummonerService } from './summoner.service';
import { LeagueApiService } from './league-api/league-api.service';
import { ClientProxy } from '@nestjs/microservices';

@Injectable()
export class SummonerSchedulerService {
  private readonly logger = new Logger(SummonerSchedulerService.name);

  constructor(
    private readonly summonerService: SummonerService,
    private readonly leagueApiService: LeagueApiService,
    @Inject('SUMMONER_SERVICE') private client: ClientProxy,
  ) {}

  @Cron('*/30 * * * * *')
  async checkInGameStatus() {
    this.logger.log('Checking in-game status for tracking summoners...');

    try {
      const trackingSummoners = await this.summonerService.findAll();

      const promises = trackingSummoners.map(async (summoner) => {
        try {
          const soloqInfo = await this.leagueApiService.getSoloQInfoByPuuid(
            summoner.puuid,
          );

          if (
            soloqInfo.wins !== summoner.wins ||
            soloqInfo.losses !== summoner.losses
          ) {
            this.logger.log(
              `Updating stats for ${summoner.gameName}#${summoner.tagLine}: wins ${summoner.wins} -> ${soloqInfo.wins}, losses ${summoner.losses} -> ${soloqInfo.losses}`,
            );

            const matchInfo = await this.leagueApiService.getLastMatchInfo(
              summoner.puuid,
            );

            await this.summonerService.update(summoner.id, {
              wins: soloqInfo.wins,
              losses: soloqInfo.losses,
              leaguePoints: soloqInfo.leaguePoints,
            });

            this.client.emit('send_summoners', {
              id: [summoner.id, soloqInfo.wins, soloqInfo.losses].join('-'),
              data: {
                summoner: `${summoner.gameName}#${summoner.tagLine}`,
                wins: soloqInfo.wins,
                losses: soloqInfo.losses,
                lpDelta: soloqInfo.leaguePoints - summoner.leaguePoints,
                leaguePoints: soloqInfo.leaguePoints,
                ...matchInfo,
              },
            });
          }
        } catch (error) {
          const message =
            error instanceof Error ? error.message : 'Unknown error';
          this.logger.error(
            `Failed to check status for ${summoner.gameName}#${summoner.tagLine}: ${message}`,
          );
        }
      });

      await Promise.allSettled(promises);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unknown error';
      this.logger.error(`Error in checkInGameStatus: ${message}`);
    }
  }
}
