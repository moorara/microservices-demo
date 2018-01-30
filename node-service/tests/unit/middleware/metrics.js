/* eslint-env mocha */
const express = require('express')
const should = require('should')
const supertest = require('supertest')

const MetricsMiddleware = require('../../../middleware/metrics')

describe('MetricsMiddleware', () => {
  let middleware
  let app, request

  beforeEach(() => {
    middleware = new MetricsMiddleware()
    app = express()
    app.use(middleware.router)
  })

  describe('http', () => {
    it('should provide default metrics', done => {
      request = supertest(app)
      request.get('/metrics')
        .expect('Content-Type', /text/)
        .expect(200)
        .end((err, res) => {
          should.not.exist(err)
          done()
        })
    })
    it('should create a logger middleware with defaults', done => {
      app.use('/resources/:id', (req, res) => res.sendStatus(200))

      request = supertest(app)
      request.get('/resources/1234')
        .expect(200)
        .end((err, res) => {
          should.not.exist(err)
          request.get('/metrics')
            .expect('Content-Type', /text/)
            .expect(200)
            .end((err, res) => {
              should.not.exist(err)
              done()
            })
        })
    })
  })
})
