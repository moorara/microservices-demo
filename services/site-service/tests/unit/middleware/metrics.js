/* eslint-env mocha */
const express = require('express')
const promClient = require('prom-client')
const should = require('should')
const supertest = require('supertest')

const { normalize } = require('../../../middleware')
const MetricsMiddleware = require('../../../middleware/metrics')

describe('MetricsMiddleware', () => {
  let register
  let middleware
  let app, request

  beforeEach(() => {
    register = new promClient.Registry()
    middleware = MetricsMiddleware.create({ register })

    app = express()
    app.use(normalize())
    app.use(middleware)

    app.get('/metrics', (req, res) => {
      res.type('text')
      res.send(register.metrics())
    })
  })

  describe('create', () => {
    it('should provide http middleware metrics', done => {
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
