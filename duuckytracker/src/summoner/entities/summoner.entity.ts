import {
  Column,
  DeleteDateColumn,
  Entity,
  Index,
  PrimaryGeneratedColumn,
} from 'typeorm';

@Entity({ name: 'summoners' })
@Index(['puuid'], {
  unique: true,
  where: '"deletedAt" IS NULL',
})
export class Summoner {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  gameName: string;

  @Column()
  tagLine: string;

  @Column()
  puuid: string;

  @Column()
  wins: number;

  @Column()
  losses: number;

  @Column()
  leaguePoints: number;

  @DeleteDateColumn()
  deletedAt?: Date | null;
}
