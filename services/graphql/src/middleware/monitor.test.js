/* eslint-env mocha */
const stream = require('stream')
const express = require('express')
const winston = require('winston')
const promClient = require('prom-client')
const opentracing = require('opentracing')
const should = require('should')
const supertest = require('supertest')

const MonitorMiddleware = require('./monitor')

describe('MonitorMiddleware', () => {
  describe('constructor', () => {
    it('should create a new middleware with defaults', () => {
      const options = { tracer: new opentracing.MockTracer() }
      const middleware = new MonitorMiddleware(options)
      should.exist(middleware.router)
      should.exist(middleware.histogram)
      should.exist(middleware.summary)
      should.exist(middleware.tracer)
      middleware.spanName.should.equal('http-request')
    })
    it('should create a new middleware with provided options', () => {
      const options = {
        winston: winston.createLogger(),
        metadata: { environment: 'test' },
        register: new promClient.Registry(),
        tracer: new opentracing.MockTracer(),
        spanName: 'test-request'
      }

      const middleware = new MonitorMiddleware(options)
      should.exist(middleware.router)
      should.exist(middleware.histogram)
      should.exist(middleware.summary)
      should.equal(middleware.tracer, options.tracer)
      should.equal(middleware.spanName, options.spanName)
    })
  })

  describe('middleware', () => {
    let logs, transformStream
    let logger, register, tracer
    let middleware
    let app, request
    let reqSpan

    const verify = (histogramName, summaryName, httpVersion, method, url, statusCode, spanName) => {
      logs[0].meta.req.httpVersion.should.equal(httpVersion)
      logs[0].meta.req.method.should.equal(method)
      logs[0].meta.res.statusCode.should.equal(statusCode)
      logs[0].meta.responseTime.should.be.aboveOrEqual(0)
      should.exist(logs[0].level)
      should.exist(logs[0].message)
      should.exist(logs[0].meta.req.url)
      should.exist(logs[0].meta.req.query)
      should.exist(logs[0].meta.req.headers)
      should.exist(logs[0].meta.req.originalUrl)

      const metrics = register.getMetricsAsJSON()

      metrics[0].name.should.equal(histogramName)
      metrics[0].type.should.equal('histogram')
      for (const val of metrics[0].values) {
        val.labels.httpVersion.should.equal(httpVersion)
        val.labels.method.should.equal(method)
        val.labels.url.should.equal(url)
        val.labels.statusCode.should.equal(statusCode)
      }

      metrics[1].name.should.equal(summaryName)
      metrics[1].type.should.equal('summary')
      for (const val of metrics[1].values) {
        val.labels.httpVersion.should.equal(httpVersion)
        val.labels.method.should.equal(method)
        val.labels.url.should.equal(url)
        val.labels.statusCode.should.equal(statusCode)
      }

      const span = tracer.report().spans[0]
      span.should.eql(reqSpan)
      span.operationName().should.equal(spanName)
      span.tags()['http.version'].should.equal(httpVersion)
      span.tags()[opentracing.Tags.HTTP_METHOD].should.equal(method)
      span.tags()[opentracing.Tags.HTTP_URL].should.equal(url)
      span.tags()[opentracing.Tags.HTTP_STATUS_CODE].should.equal(statusCode)
    }

    beforeEach(() => {
      logs = []
      transformStream = new stream.PassThrough()
      transformStream.on('data', data => {
        logs.push(JSON.parse(data))
      })

      logger = winston.createLogger({
        transports: [
          new winston.transports.Stream({ stream: transformStream })
        ]
      })

      register = new promClient.Registry()
      tracer = new opentracing.MockTracer()
      middleware = new MonitorMiddleware({ winston: logger, register, tracer })

      app = express()
      app.use(middleware.router)
      app.use((req, res, next) => {
        reqSpan = req.context.span
        res.sendStatus(200)
      })

      request = supertest(app)
    })

    it('should generate logs, metrics, and traces', done => {
      request.get('/graphql')
        .expect(200)
        .then(res => {
          verify(
            'http_requests_duration_seconds',
            'http_requests_duration_quantiles_seconds',
            '1.1', 'GET', '/graphql', 200,
            'http-request'
          )
          done()
        }).catch(done)
    })
    it('should generate logs, metrics, and traces', done => {
      request.get('/v1/graphql?query=someting')
        .expect(200)
        .then(res => {
          verify(
            'http_requests_duration_seconds',
            'http_requests_duration_quantiles_seconds',
            '1.1', 'GET', '/v1/graphql', 200,
            'http-request'
          )
          done()
        }).catch(done)
    })
  })
})
