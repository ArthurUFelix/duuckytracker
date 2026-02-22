import { Test, TestingModule } from '@nestjs/testing';
import { SummonerController } from './summoner.controller';
import { SummonerService } from './summoner.service';
import { CreateSummonerDto } from './dto/create-summoner.dto';

describe('SummonerController', () => {
  let controller: SummonerController;

  const mockSummonerService = {
    create: jest.fn(),
    findAll: jest.fn(),
    remove: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [SummonerController],
      providers: [
        {
          provide: SummonerService,
          useValue: mockSummonerService,
        },
      ],
    }).compile();

    controller = module.get<SummonerController>(SummonerController);
  });

  describe('create', () => {
    it('should create a summoner', async () => {
      const createSummonerDto: CreateSummonerDto = {
        gameName: 'TestSummoner',
        tagLine: 'NA1',
      };
      const result = {
        id: 1,
        ...createSummonerDto,
        puuid: 'test-puuid',
        wins: 150,
        losses: 120,
        leaguePoints: 75,
        deletedAt: null,
      };
      mockSummonerService.create.mockResolvedValue(result);

      expect(await controller.create(createSummonerDto)).toBe(result);
      expect(mockSummonerService.create).toHaveBeenCalledWith(
        createSummonerDto,
      );
    });
  });

  describe('findAll', () => {
    it('should return all summoners', async () => {
      const result = [
        {
          id: 1,
          gameName: 'Summoner1',
          tagLine: 'NA1',
          puuid: 'puuid1',
          wins: 150,
          losses: 120,
          leaguePoints: 75,
          deletedAt: null,
        },
        {
          id: 2,
          gameName: 'Summoner2',
          tagLine: 'EUW',
          puuid: 'puuid2',
          wins: 200,
          losses: 180,
          leaguePoints: 50,
          deletedAt: null,
        },
      ];
      mockSummonerService.findAll.mockResolvedValue(result);

      expect(await controller.findAll()).toBe(result);
      expect(mockSummonerService.findAll).toHaveBeenCalled();
    });
  });

  describe('remove', () => {
    it('should remove a summoner by id', async () => {
      const id = 1;
      mockSummonerService.remove.mockResolvedValue(undefined);

      expect(await controller.remove(id)).toBeUndefined();
      expect(mockSummonerService.remove).toHaveBeenCalledWith(id);
    });
  });
});
