import { IsNotEmpty, IsString } from 'class-validator';

export class CreateSummonerDto {
  @IsString()
  @IsNotEmpty()
  gameName: string;

  @IsString()
  @IsNotEmpty()
  tagLine: string;
}
