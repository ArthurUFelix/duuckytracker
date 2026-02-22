import { Test, TestingModule } from '@nestjs/testing';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { HttpException, HttpStatus } from '@nestjs/common';
import { SummonerService } from './summoner.service';
import { Summoner } from './entities/summoner.entity';
import { CreateSummonerDto } from './dto/create-summoner.dto';
import { LeagueApiService } from './league-api/league-api.service';

describe('SummonerService', () => {
  let service: SummonerService;
  let summonerRepository: Repository<Summoner>;
  let leagueApiService: LeagueApiService;

  const mockSummoner: Summoner = {
    id: 1,
    gameName: 'TestSummoner',
    tagLine: 'NA1',
    puuid: 'test-puuid-123',
    wins: 150,
    losses: 120,
    leaguePoints: 75,
    deletedAt: null,
  };

  const mockCreateSummonerDto: CreateSummonerDto = {
    gameName: 'TestSummoner',
    tagLine: 'NA1',
  };

  const mockApiResponse = {
    puuid: 'test-puuid-123',
    gameName: 'TestSummoner',
    tagLine: 'NA1',
  };

  const mockLeagueApiService = {
    getSummonerByName: jest.fn(),
    getSoloQInfoByPuuid: jest.fn(),
  };

  let mockRepository: any;

  beforeEach(async () => {
    mockRepository = {
      create: jest.fn(),
      save: jest.fn(),
      find: jest.fn(),
      softDelete: jest.fn(),
      update: jest.fn(),
      findOne: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        SummonerService,
        {
          provide: getRepositoryToken(Summoner),
          useValue: mockRepository,
        },
        {
          provide: LeagueApiService,
          useValue: mockLeagueApiService,
        },
      ],
    }).compile();

    service = module.get<SummonerService>(SummonerService);
    summonerRepository = module.get<Repository<Summoner>>(
      getRepositoryToken(Summoner),
    );
    leagueApiService = module.get<LeagueApiService>(LeagueApiService);

    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should successfully create a summoner', async () => {
      const mockSoloQInfo = {
        wins: 150,
        losses: 120,
        leaguePoints: 75,
      };
      mockLeagueApiService.getSummonerByName.mockResolvedValue(mockApiResponse);
      mockLeagueApiService.getSoloQInfoByPuuid.mockResolvedValue(mockSoloQInfo);
      mockRepository.create.mockReturnValue(mockSummoner);
      mockRepository.save.mockResolvedValue(mockSummoner);

      const result = await service.create(mockCreateSummonerDto);

      expect(mockLeagueApiService.getSummonerByName).toHaveBeenCalledWith(
        'TestSummoner',
        'NA1',
      );
      expect(mockLeagueApiService.getSoloQInfoByPuuid).toHaveBeenCalledWith(
        'test-puuid-123',
      );
      expect(mockRepository.create).toHaveBeenCalledWith({
        ...mockCreateSummonerDto,
        puuid: mockApiResponse.puuid,
        wins: mockSoloQInfo.wins,
        losses: mockSoloQInfo.losses,
        leaguePoints: mockSoloQInfo.leaguePoints,
      });
      expect(mockRepository.save).toHaveBeenCalledWith(mockSummoner);
      expect(result).toEqual(mockSummoner);
    });

    it('should throw HttpException when summoner already exists (duplicate key error)', async () => {
      const mockSoloQInfo = {
        wins: 150,
        losses: 120,
        leaguePoints: 75,
      };
      mockLeagueApiService.getSummonerByName.mockResolvedValue(mockApiResponse);
      mockLeagueApiService.getSoloQInfoByPuuid.mockResolvedValue(mockSoloQInfo);
      mockRepository.create.mockReturnValue(mockSummoner);
      const duplicateError = new Error(
        'duplicate key value violates unique constraint',
      );
      duplicateError['code'] = '23505';
      mockRepository.save.mockRejectedValue(duplicateError);

      await expect(service.create(mockCreateSummonerDto)).rejects.toThrow(
        new HttpException('Summoner already exists', HttpStatus.BAD_REQUEST),
      );
      expect(mockLeagueApiService.getSummonerByName).toHaveBeenCalled();
      expect(mockLeagueApiService.getSoloQInfoByPuuid).toHaveBeenCalled();
      expect(mockRepository.create).toHaveBeenCalled();
      expect(mockRepository.save).toHaveBeenCalled();
    });
  });

  describe('findAll', () => {
    it('should return all summoners', async () => {
      const mockSummoners = [mockSummoner];
      mockRepository.find.mockResolvedValue(mockSummoners);

      const result = await service.findAll();

      expect(mockRepository.find).toHaveBeenCalled();
      expect(result).toEqual(mockSummoners);
    });

    it('should return empty array when no summoners exist', async () => {
      mockRepository.find.mockResolvedValue([]);

      const result = await service.findAll();

      expect(mockRepository.find).toHaveBeenCalled();
      expect(result).toEqual([]);
    });
  });

  describe('remove', () => {
    it('should successfully remove a summoner by id', async () => {
      const deleteResult = { affected: 1, raw: [] };
      mockRepository.softDelete.mockResolvedValue(deleteResult);

      const result = await service.remove(1);

      expect(mockRepository.softDelete).toHaveBeenCalledWith(1);
      expect(result).toEqual(deleteResult);
    });

    it('should handle removing non-existent summoner', async () => {
      const deleteResult = { affected: 0, raw: [] };
      mockRepository.softDelete.mockResolvedValue(deleteResult);

      const result = await service.remove(999);

      expect(mockRepository.softDelete).toHaveBeenCalledWith(999);
      expect(result).toEqual(deleteResult);
    });
  });

  describe('update', () => {
    it('should successfully update a summoner', async () => {
      const updateData = { wins: 200, losses: 150 };
      const updatedSummoner = { ...mockSummoner, ...updateData };
      mockRepository.update.mockResolvedValue({ affected: 1, raw: [] });
      mockRepository.findOne.mockResolvedValue(updatedSummoner);

      const result = await service.update(1, updateData);

      expect(mockRepository.update).toHaveBeenCalledWith(1, updateData);
      expect(mockRepository.findOne).toHaveBeenCalledWith({ where: { id: 1 } });
      expect(result).toEqual(updatedSummoner);
    });

    it('should return null if summoner not found', async () => {
      mockRepository.update.mockResolvedValue({ affected: 0, raw: [] });
      mockRepository.findOne.mockResolvedValue(null);

      const result = await service.update(999, { wins: 100 });

      expect(mockRepository.update).toHaveBeenCalledWith(999, { wins: 100 });
      expect(result).toBeNull();
    });
  });
});
