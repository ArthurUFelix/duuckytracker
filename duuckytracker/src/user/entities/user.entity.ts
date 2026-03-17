import {
  Column,
  DeleteDateColumn,
  Entity,
  Index,
  PrimaryGeneratedColumn,
} from 'typeorm';

@Entity({ name: 'users' })
@Index(['username'], {
  unique: true,
  where: '"deletedAt" IS NULL',
})
export class User {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  username: string;

  @Column()
  password: string;

  @Column('text', { array: true, default: () => 'ARRAY[]::text[]' })
  trackedSummoners: string[];

  @DeleteDateColumn()
  deletedAt?: Date | null;
}
