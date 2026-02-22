import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { App } from 'supertest/types';
import { AppModule } from './../src/app.module';

describe('SummonerController (e2e)', () => {
  let app: INestApplication<App>;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    await app.init();
  });

  afterEach(async () => {
    await app.close();
  });

  describe('/summoner (GET)', () => {
    it('should return an array of summoners', () => {
      return request(app.getHttpServer())
        .get('/summoner')
        .expect(200)
        .expect([]);
    });
  });

  describe('/summoner (POST)', () => {
    it('should create a summoner', () => {
      const createSummonerDto = {
        gameName: 'TestSummoner',
        tagLine: 'NA1',
      };

      return request(app.getHttpServer())
        .post('/summoner')
        .send(createSummonerDto)
        .expect(201);
    });
  });

  describe('/summoner/:id (DELETE)', () => {
    it('should delete a summoner', () => {
      return request(app.getHttpServer()).delete('/summoner/1').expect(200);
    });
  });
});
