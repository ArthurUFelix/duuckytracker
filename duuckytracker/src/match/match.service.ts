import { Injectable } from '@nestjs/common';
import { CreateMatchDto } from './dto/create-match.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Match } from './entities/match.entity';

@Injectable()
export class MatchService {
  constructor(
    @InjectRepository(Match)
    private matchRepository: Repository<Match>,
  ) {}

  async create(createMathDto: CreateMatchDto) {
    const newMatch = this.matchRepository.create(createMathDto);
    return await this.matchRepository.save(newMatch);
  }

  async findAll(filters = {}) {
    return await this.matchRepository.find(filters);
  }
}
