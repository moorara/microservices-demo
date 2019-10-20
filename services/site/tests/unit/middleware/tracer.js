/* eslint-env mocha */
require('should')
const sinon = require('sinon')
const supertest = require('supertest')
const express = require('express')
const opentracing = require('opentracing')

const { normalize } = require('../../../middleware')
const TracerMiddleware = require('../../../middleware/tracer')

describe('TracerMiddleware', () => {
  describe('create', () => {
    let tracer, _tracer
    let middleware
    let app, request
    let spanFromCtx

    const verify = (spanName, httpVersion, method, url, statusCode) => {
      const span = tracer.report().spans[0]
      span.should.eql(spanFromCtx)
      span.operationName().should.equal(spanName)
      span.tags()['http.version'].should.equal(httpVersion)
      span.tags()[opentracing.Tags.HTTP_METHOD].should.equal(method)
      span.tags()[opentracing.Tags.HTTP_URL].should.equal(url)
      span.tags()[opentracing.Tags.HTTP_STATUS_CODE].should.equal(statusCode)
    }

    beforeEach(() => {
      // MockTracer has not implemented inect and extract!
      tracer = new opentracing.MockTracer()
      _tracer = sinon.mock(tracer)

      middleware = TracerMiddleware.create({ tracer })

      app = express()
      app.use(normalize())
      app.use(middleware)
      app.use('/resources/:id', (req, res) => {
        spanFromCtx = req.context.span
        res.sendStatus(200)
      })
    })

    afterEach(() => {
      _tracer.restore()
    })

    it('should create a span and add it to request context', done => {
      _tracer.expects('extract').withArgs(opentracing.FORMAT_HTTP_HEADERS).returns()

      request = supertest(app)
      request.get('/resources/abcd')
        .expect(200)
        .then(res => {
          _tracer.verify()
          verify('http-request', '1.1', 'GET', '/resources/abcd', 200)
          done()
        }).catch(done)
    })
    it('should create a span without parent context and add it to request context', done => {
      _tracer.expects('extract').withArgs(opentracing.FORMAT_HTTP_HEADERS).returns()

      request = supertest(app)
      request.get('/resources/1234')
        .expect(200)
        .then(res => {
          _tracer.verify()
          verify('http-request', '1.1', 'GET', '/resources/1234', 200)
          done()
        }).catch(done)
    })
    it('should create a span with parent context and add it to request context', done => {
      const parentSpanContext = Object.assign(new opentracing.SpanContext(), {
        _traceIdStr: 'd0d03caa2f19335f',
        _spanIdStr: '85abb56a99573b64',
        _parentIdStr: 'd0d03caa2f19335f',
        _flags: 1,
        _baggage: {}
      })
      _tracer.expects('extract').withArgs(opentracing.FORMAT_HTTP_HEADERS).returns(parentSpanContext)

      request = supertest(app)
      request.get('/resources/1234')
        .set('uber-trace-id', 'd0d03caa2f19335f:85abb56a99573b64:d0d03caa2f19335f:1')
        .expect(200)
        .then(res => {
          _tracer.verify()
          verify('http-request', '1.1', 'GET', '/resources/1234', 200)
          done()
        }).catch(done)
    })
  })
})
