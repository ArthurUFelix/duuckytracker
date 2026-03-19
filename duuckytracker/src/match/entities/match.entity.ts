import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity({ name: 'matches' })
export class Match {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  championId: string;

  @Column()
  summonerName: string;

  @Column()
  date: Date;
}
