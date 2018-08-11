/* eslint-env mocha */
require('should')
const supertest = require('supertest')
const express = require('express')
const opentracing = require('opentracing')

const { normalize } = require('../../../middleware')
const TracerMiddleware = require('../../../middleware/tracer')

describe('TracerMiddleware', () => {
  describe('create', () => {
    let tracer, middleware
    let app, request
    let spanFromCtx

    const verify = (spanName, method, url, statusCode) => {
      const span = tracer._spans[0]
      span.should.eql(spanFromCtx)
      span._operationName.should.equal(spanName)
      span._tags['http.method'].should.equal(method)
      span._tags['http.url'].should.equal(url)
      span._tags['http.status_code'].should.equal(statusCode)
    }

    beforeEach(() => {
      tracer = new opentracing.MockTracer()
      middleware = TracerMiddleware.create({ tracer })

      app = express()
      app.use(normalize())
      app.use(middleware)
      app.use('/resources/:id', (req, res) => {
        spanFromCtx = req.context.span
        res.sendStatus(200)
      })
    })

    it('should create a span and add it to request context', done => {
      request = supertest(app)
      request.get('/resources/1234')
        .expect(200)
        .then(res => {
          verify('http-request', 'GET', '/resources/1234', 200)
          done()
        }).catch(done)
    })
    it('should create a span and add it to request context', done => {
      request = supertest(app)
      request.get('/resources/abcd')
        .expect(200)
        .then(res => {
          verify('http-request', 'GET', '/resources/abcd', 200)
          done()
        }).catch(done)
    })
  })
})
