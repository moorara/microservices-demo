/* eslint-env mocha */
const express = require('express')
const promClient = require('prom-client')
const should = require('should')
const supertest = require('supertest')

const Metrics = require('../../../util/metrics')

describe('Metrics', () => {
  describe('constructor', () => {
    let metrics

    it('should create new metrics with global registry', () => {
      metrics = new Metrics()
      should.exist(metrics.register)
      should.exist(metrics.router)
    })
    it('should create new metrics with a new registry', () => {
      const register = new promClient.Registry()
      metrics = new Metrics({ register })
      should.exist(metrics.register)
      should.exist(metrics.router)
    })
  })

  describe('getMetrics', () => {
    let metrics
    let app

    beforeEach(() => {
      metrics = new Metrics({
        register: new promClient.Registry()
      })
      app = express()
      app.use(metrics.router)
    })

    it('should respond with default metrics', done => {
      const request = supertest(app)
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
