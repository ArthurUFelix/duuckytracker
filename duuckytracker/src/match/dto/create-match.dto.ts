import { IsDateString, IsNotEmpty, IsString } from 'class-validator';

export class CreateMatchDto {
  @IsString()
  @IsNotEmpty()
  champion: string;

  @IsDateString()
  @IsNotEmpty()
  date: string;
}
