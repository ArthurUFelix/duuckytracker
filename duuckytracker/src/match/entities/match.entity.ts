import {
  Column,
  DeleteDateColumn,
  Entity,
  Index,
  PrimaryGeneratedColumn,
} from 'typeorm';

@Entity({ name: 'matches' })
export class Match {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  champion: string;

  @Column()
  date: Date;
}
