import {
  Controller,
  Get,
  Post,
  Body,
  Param,
  Delete,
  ParseIntPipe,
} from '@nestjs/common';
import { SummonerService } from './summoner.service';
import { CreateSummonerDto } from './dto/create-summoner.dto';

@Controller('summoner')
export class SummonerController {
  constructor(private readonly summonerService: SummonerService) {}

  @Post()
  create(@Body() createSummonerDto: CreateSummonerDto) {
    return this.summonerService.create(createSummonerDto);
  }

  @Get()
  findAll() {
    return this.summonerService.findAll();
  }

  @Delete(':id')
  remove(@Param('id', ParseIntPipe) id: number) {
    return this.summonerService.remove(id);
  }
}
