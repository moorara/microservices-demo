/* eslint-env mocha */
require('should')
const supertest = require('supertest')
const express = require('express')
const opentracing = require('opentracing')

const { normalize } = require('../../../middleware')
const TracerMiddleware = require('../../../middleware/tracer')

describe('TracerMiddleware', () => {
  let middleware
  let app, request

  describe('http', () => {
    it('should create a tracer middleware', done => {
      const tracer = new opentracing.MockTracer()
      middleware = TracerMiddleware.http({ tracer })

      app = express()
      app.use(normalize())
      app.use(middleware)
      app.use('/resources/:id', (req, res) => res.sendStatus(200))

      request = supertest(app)
      request.get('/resources/1234')
        .expect(200)
        .then(res => {
          tracer._spans.should.have.length(1)
          tracer._spans[0]._tags['http.method'].should.equal('GET')
          tracer._spans[0]._tags['http.url'].should.equal('/resources/:id')
          tracer._spans[0]._tags['http.status_code'].should.equal(200)
          done()
        })
    })
  })
})
