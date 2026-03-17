import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { CreateSummonerDto } from './dto/create-summoner.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Summoner } from './entities/summoner.entity';
import { Repository } from 'typeorm';
import { LeagueApiService } from './league-api/league-api.service';

@Injectable()
export class SummonerService {
  constructor(
    @InjectRepository(Summoner)
    private summonerRepository: Repository<Summoner>,
    private leagueApiService: LeagueApiService,
  ) {}

  async create(createSummonerDto: CreateSummonerDto) {
    try {
      const summonerData = await this.leagueApiService.getSummonerByName(
        createSummonerDto.gameName,
        createSummonerDto.tagLine,
      );

      const soloqInfo = await this.leagueApiService.getSoloQInfoByPuuid(
        summonerData.puuid,
      );

      const newSummoner = this.summonerRepository.create({
        ...createSummonerDto,
        puuid: summonerData.puuid,
        wins: soloqInfo.wins,
        losses: soloqInfo.losses,
        leaguePoints: soloqInfo.leaguePoints,
      });

      return await this.summonerRepository.save(newSummoner);
    } catch (error) {
      if (error.code === '23505') {
        throw new HttpException(
          'Summoner already exists',
          HttpStatus.BAD_REQUEST,
        );
      }
      throw error;
    }
  }

  async findAll() {
    return await this.summonerRepository.find();
  }

  async remove(id: number) {
    return await this.summonerRepository.softDelete(id);
  }

  async update(id: number, updateData: Partial<Summoner>) {
    await this.summonerRepository.update(id, updateData);
    return await this.summonerRepository.findOne({ where: { id } });
  }
}
