import { Injectable, HttpException, HttpStatus } from '@nestjs/common';

@Injectable()
export class LeagueApiService {
  private readonly accountAPI =
    'https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/';
  private readonly queueAPI =
    'https://br1.api.riotgames.com/lol/league/v4/entries/by-puuid/';
  private readonly apiKey = process.env.RIOT_API_KEY;

  async getSummonerByName(gameName: string, tagLine: string) {
    if (!this.apiKey) {
      throw new HttpException(
        'Riot API key not configured',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }

    try {
      const response = await fetch(
        `${this.accountAPI}${encodeURIComponent(gameName)}/${encodeURIComponent(tagLine)}`,
        {
          headers: {
            'X-Riot-Token': this.apiKey,
          },
        },
      );

      if (!response.ok) {
        if (response.status === 404) {
          throw new HttpException('Summoner not found', HttpStatus.NOT_FOUND);
        } else {
          throw new HttpException(
            `Riot API error: ${response.status}`,
            HttpStatus.BAD_GATEWAY,
          );
        }
      }

      const summonerData = await response.json();
      return summonerData;
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }
      throw new HttpException(
        'Failed to fetch summoner data',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }

  async getSoloQInfoByPuuid(puuid: string) {
    if (!this.apiKey) {
      throw new HttpException(
        'Riot API key not configured',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }

    try {
      const response = await fetch(
        `${this.queueAPI}${encodeURIComponent(puuid)}`,
        {
          headers: {
            'X-Riot-Token': this.apiKey,
          },
        },
      );

      if (!response.ok) {
        throw new HttpException(
          `Riot API error: ${response.status}`,
          HttpStatus.BAD_GATEWAY,
        );
      }

      const queueInfo = await response.json();
      const soloq = queueInfo.find(
        (entry) => entry.queueType === 'RANKED_SOLO_5x5',
      );
      return soloq || { wins: 0, losses: 0, leaguePoints: 0 };
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }
      throw new HttpException(
        'Failed to fetch summoner data',
        HttpStatus.INTERNAL_SERVER_ERROR,
      );
    }
  }
}
