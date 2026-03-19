import { IsDateString, IsNotEmpty, IsString } from 'class-validator';

export class CreateMatchDto {
  @IsString()
  @IsNotEmpty()
  championId: string;

  @IsString()
  @IsNotEmpty()
  summonerName: string;

  @IsDateString()
  @IsNotEmpty()
  date: string;
}
