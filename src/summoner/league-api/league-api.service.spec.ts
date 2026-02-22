import { Test, TestingModule } from '@nestjs/testing';
import { HttpException, HttpStatus } from '@nestjs/common';
import { LeagueApiService } from './league-api.service';

const mockFetch = jest.fn();
global.fetch = mockFetch;

describe('LeagueApiService', () => {
  let service: LeagueApiService;

  beforeEach(async () => {
    jest.clearAllMocks();
    process.env.RIOT_API_KEY = 'test-api-key';

    const module: TestingModule = await Test.createTestingModule({
      providers: [LeagueApiService],
    }).compile();
    service = module.get<LeagueApiService>(LeagueApiService);
  });

  afterEach(() => {
    delete process.env.RIOT_API_KEY;
  });

  describe('getSummonerByName', () => {
    const mockSummonerData = {
      puuid: 'test-puuid-123',
      gameName: 'TestSummoner',
      tagLine: 'NA1',
    };

    it('should return summoner data on successful API call', async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockSummonerData),
      };
      mockFetch.mockResolvedValue(mockResponse);

      const result = await service.getSummonerByName('TestSummoner', 'NA1');

      expect(mockFetch).toHaveBeenCalledWith(
        'https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/TestSummoner/NA1',
        {
          headers: {
            'X-Riot-Token': 'test-api-key',
          },
        },
      );
      expect(result).toEqual(mockSummonerData);
    });

    it('should throw error when API key is not configured', async () => {
      delete process.env.RIOT_API_KEY;

      const module: TestingModule = await Test.createTestingModule({
        providers: [LeagueApiService],
      }).compile();

      service = module.get<LeagueApiService>(LeagueApiService);

      await expect(
        service.getSummonerByName('TestSummoner', 'NA1'),
      ).rejects.toThrow(
        new HttpException(
          'Riot API key not configured',
          HttpStatus.INTERNAL_SERVER_ERROR,
        ),
      );
    });

    it('should throw 404 error when summoner is not found', async () => {
      const mockResponse = {
        ok: false,
        status: 404,
      };
      mockFetch.mockResolvedValue(mockResponse);

      await expect(
        service.getSummonerByName('NonExistent', 'NA1'),
      ).rejects.toThrow(
        new HttpException('Summoner not found', HttpStatus.NOT_FOUND),
      );
    });

    it('should throw generic error for other HTTP status codes', async () => {
      const mockResponse = {
        ok: false,
        status: 500,
      };
      mockFetch.mockResolvedValue(mockResponse);

      await expect(
        service.getSummonerByName('TestSummoner', 'NA1'),
      ).rejects.toThrow(
        new HttpException('Riot API error: 500', HttpStatus.BAD_GATEWAY),
      );
    });
  });

  describe('getSoloQInfoByPuuid', () => {
    const mockQueueData = [
      {
        queueType: 'RANKED_SOLO_5x5',
        wins: 150,
        losses: 120,
        leaguePoints: 75,
      },
      {
        queueType: 'RANKED_FLEX_SR',
        wins: 50,
        losses: 40,
        leaguePoints: 25,
      },
    ];

    it('should return solo queue info for existing summoner', async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue(mockQueueData),
      };
      mockFetch.mockResolvedValue(mockResponse);

      const result = await service.getSoloQInfoByPuuid('test-puuid');

      expect(mockFetch).toHaveBeenCalledWith(
        'https://br1.api.riotgames.com/lol/league/v4/entries/by-puuid/test-puuid',
        {
          headers: {
            'X-Riot-Token': 'test-api-key',
          },
        },
      );
      expect(result).toEqual({
        queueType: 'RANKED_SOLO_5x5',
        wins: 150,
        losses: 120,
        leaguePoints: 75,
      });
    });

    it('should return default values when no solo queue data exists', async () => {
      const mockResponse = {
        ok: true,
        json: jest.fn().mockResolvedValue([]),
      };
      mockFetch.mockResolvedValue(mockResponse);

      const result = await service.getSoloQInfoByPuuid('test-puuid');

      expect(result).toEqual({
        wins: 0,
        losses: 0,
        leaguePoints: 0,
      });
    });

    it('should throw error when API key is not configured', async () => {
      delete process.env.RIOT_API_KEY;

      const module: TestingModule = await Test.createTestingModule({
        providers: [LeagueApiService],
      }).compile();

      service = module.get<LeagueApiService>(LeagueApiService);

      await expect(service.getSoloQInfoByPuuid('test-puuid')).rejects.toThrow(
        new HttpException(
          'Riot API key not configured',
          HttpStatus.INTERNAL_SERVER_ERROR,
        ),
      );
    });

    it('should throw generic error for HTTP failure', async () => {
      const mockResponse = {
        ok: false,
        status: 500,
      };
      mockFetch.mockResolvedValue(mockResponse);

      await expect(service.getSoloQInfoByPuuid('test-puuid')).rejects.toThrow(
        new HttpException('Riot API error: 500', HttpStatus.BAD_GATEWAY),
      );
    });
  });
});
