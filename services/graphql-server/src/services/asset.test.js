/* eslint-env mocha */
const { URL } = require('url')
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const AssetService = require('./asset')

describe('AssetService', () => {
  let logger
  let histogram, _histogram
  let summary, _summary
  let tracer, _tracer
  let nats, _nats
  let config, options
  let service, _service
  let span, context

  beforeEach(() => {
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    histogram = {
      observe () {}
    }
    _histogram = sinon.mock(histogram)

    summary = {
      observe () {}
    }
    _summary = sinon.mock(summary)

    // MockTracer has not implemented inect and extract!
    tracer = new opentracing.MockTracer()
    _tracer = sinon.mock(tracer)
    _tracer.expects('inject').returns()

    nats = {
      publish () {},
      subscribe () {},
      request () {},
      requestOne () {},
      currentServer: {
        url: new URL('nats://nats:4222')
      }
    }
    _nats = sinon.mock(nats)

    config = {
      servers: 'nats://nats:4222',
      user: 'client',
      pass: 'pass'
    }
    options = { logger, histogram, summary, tracer, nats }

    service = new AssetService(config, options)
    _service = sinon.mock(service)

    span = {}
    context = { span }
  })

  afterEach(() => {
    _histogram.restore()
    _summary.restore()
    _tracer.restore()
    _nats.restore()
    _service.restore()
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new AssetService(config, { tracer: options.tracer, nats: options.nats })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
      should.exist(service.nats)
    })
    it('should create a new service with provided options', () => {
      const service = new AssetService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
      service.nats.should.equal(options.nats)
    })
  })

  describe('exec', () => {
    const verifyTrace = (spanName, logMessage) => {
      const span = tracer.report().spans[0]
      span.operationName().should.equal(spanName)
      span.tags()[opentracing.Tags.SPAN_KIND].should.equal('requester')
      span.tags()[opentracing.Tags.PEER_SERVICE].should.equal('asset-service')
      span.tags()[opentracing.Tags.PEER_ADDRESS].should.equal(nats.currentServer.url.host)
      span._logs[0].fields.event.should.equal(spanName)
      span._logs[0].fields.message.should.equal(logMessage)
    }

    it('should reject with an error when input function rejects', done => {
      const stub = sinon.stub()
      stub.rejects(new Error('mock error'))
      _histogram.expects('observe').withArgs({ op: 'unit-test', success: 'false' }).returns()
      _summary.expects('observe').withArgs({ op: 'unit-test', success: 'false' }).returns()
      service.exec(context, 'unit-test', stub).catch(err => {
        err.message.should.equal('mock error')
        _histogram.verify()
        _summary.verify()
        verifyTrace('unit-test', 'mock error')
        done()
      })
    })
    it('should resolve successfully when input function resolves', done => {
      const stub = sinon.stub()
      stub.resolves()
      _histogram.expects('observe').withArgs({ op: 'unit-test', success: 'true' }).returns()
      _summary.expects('observe').withArgs({ op: 'unit-test', success: 'true' }).returns()
      service.exec(context, 'unit-test', stub).then(result => {
        _histogram.verify()
        _summary.verify()
        verifyTrace('unit-test', 'successful!')
        done()
      }).catch(done)
    })
  })

  describe('getAsset', () => {

  })

  describe('getAssets', () => {

  })

  describe('createAlarm', () => {

  })

  describe('updateAlarm', () => {

  })

  describe('deleteAlarm', () => {

  })

  describe('createCamera', () => {

  })

  describe('updateCamera', () => {

  })

  describe('deleteCamera', () => {

  })
})
