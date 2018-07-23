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
          tracer._spans[0].should.eql(spanFromCtx)
          tracer._spans[0]._operationName.should.equal('http-request')
          tracer._spans[0]._tags['http.method'].should.equal('GET')
          tracer._spans[0]._tags['http.url'].should.equal('/resources/1234')
          tracer._spans[0]._tags['http.status_code'].should.equal(200)
          done()
        })
    })
  })
})
