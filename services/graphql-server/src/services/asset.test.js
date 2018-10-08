/* eslint-env mocha */
const { URL } = require('url')
const nats = require('nats')
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const AssetService = require('./asset')

const expectedSubject = 'asset_service'
const expectedOptions = {}
const expectedTimeout = 1000

describe('AssetService', () => {
  let logger
  let histogram, _histogram
  let summary, _summary
  let tracer, _tracer
  let conn, _conn
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

    conn = {
      publish () {},
      subscribe () {},
      request () {},
      requestOne () {},
      currentServer: {
        url: new URL('nats://nats:4222')
      }
    }
    _conn = sinon.mock(conn)

    config = {
      servers: 'nats://nats:4222',
      user: 'client',
      pass: 'pass'
    }
    options = { logger, histogram, summary, tracer, conn }

    service = new AssetService(config, options)
    _service = sinon.mock(service)

    span = {}
    context = { span }
  })

  afterEach(() => {
    _histogram.restore()
    _summary.restore()
    _tracer.restore()
    _conn.restore()
    _service.restore()
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new AssetService(config, { tracer: options.tracer, conn: options.conn })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
      should.exist(service.conn)
    })
    it('should create a new service with provided options', () => {
      const service = new AssetService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
      service.conn.should.equal(options.conn)
    })
  })

  describe('exec', () => {
    const verifyTrace = (spanName, logMessage) => {
      const span = tracer.report().spans[0]
      span.operationName().should.equal(spanName)
      span.tags()[opentracing.Tags.SPAN_KIND].should.equal('requester')
      span.tags()[opentracing.Tags.PEER_SERVICE].should.equal('asset-service')
      span.tags()[opentracing.Tags.PEER_ADDRESS].should.equal(conn.currentServer.url.host)
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

  describe('createAlarm', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('create-alarm')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const requestMsg = JSON.stringify({ kind: 'createAlarm', span: '{}', input })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.createAlarm(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const requestMsg = JSON.stringify({ kind: 'createAlarm', span: '{}', input })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createAlarm(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const requestMsg = JSON.stringify({ kind: 'createAlarm', span: '{}', input })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createAlarm(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const alarm = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const requestMsg = JSON.stringify({ kind: 'createAlarm', span: '{}', input })
      const responseMsg = JSON.stringify({ kind: 'createAlarm', alarm })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createAlarm(context, input).then(result => {
        result.should.eql(alarm)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('allAlarm', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('all-alarm')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allAlarm', span: '{}', siteId })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.allAlarm(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allAlarm', span: '{}', siteId })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allAlarm(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allAlarm', span: '{}', siteId })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allAlarm(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const siteId = '1111-1111'
      const alarms = [{ id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'co' }]
      const requestMsg = JSON.stringify({ kind: 'allAlarm', span: '{}', siteId })
      const responseMsg = JSON.stringify({ kind: 'allAlarm', alarms })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allAlarm(context, siteId).then(result => {
        result.should.eql(alarms)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('getAlarm', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-alarm')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'getAlarm', span: '{}', id })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.getAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'getAlarm', span: '{}', id })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'getAlarm', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'aaaa-aaaa'
      const alarm = { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'co' }
      const requestMsg = JSON.stringify({ kind: 'getAlarm', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'getAlarm', alarm })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getAlarm(context, id).then(result => {
        result.should.eql(alarm)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('updateAlarm', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('update-alarm')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'aaaa-aaaa'
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
      const requestMsg = JSON.stringify({ kind: 'updateAlarm', span: '{}', id, input })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.updateAlarm(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'aaaa-aaaa'
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
      const requestMsg = JSON.stringify({ kind: 'updateAlarm', span: '{}', id, input })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateAlarm(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'aaaa-aaaa'
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
      const requestMsg = JSON.stringify({ kind: 'updateAlarm', span: '{}', id, input })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateAlarm(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'aaaa-aaaa'
      const input = { siteId: '1111-1111', serialNo: '1001', material: 'smoke' }
      const requestMsg = JSON.stringify({ kind: 'updateAlarm', span: '{}', id, input })
      const responseMsg = JSON.stringify({ kind: 'updateAlarm', updated: true })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateAlarm(context, id, input).then(result => {
        result.should.be.true()
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('deleteAlarm', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('delete-alarm')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'deleteAlarm', span: '{}', id })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.deleteAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'deleteAlarm', span: '{}', id })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'deleteAlarm', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteAlarm(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'aaaa-aaaa'
      const requestMsg = JSON.stringify({ kind: 'deleteAlarm', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'deleteAlarm', deleted: true })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteAlarm(context, id).then(result => {
        result.should.be.true()
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('createCamera', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('create-camera')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const requestMsg = JSON.stringify({ kind: 'createCamera', span: '{}', input })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.createCamera(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const requestMsg = JSON.stringify({ kind: 'createCamera', span: '{}', input })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createCamera(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const requestMsg = JSON.stringify({ kind: 'createCamera', span: '{}', input })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createCamera(context, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const camera = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const requestMsg = JSON.stringify({ kind: 'createCamera', span: '{}', input })
      const responseMsg = JSON.stringify({ kind: 'createCamera', camera })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.createCamera(context, input).then(result => {
        result.should.eql(camera)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('allCamera', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('all-camera')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allCamera', span: '{}', siteId })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.allCamera(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allCamera', span: '{}', siteId })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allCamera(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const siteId = '1111-1111'
      const requestMsg = JSON.stringify({ kind: 'allCamera', span: '{}', siteId })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allCamera(context, siteId).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const siteId = '1111-1111'
      const cameras = [{ id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }]
      const requestMsg = JSON.stringify({ kind: 'allCamera', span: '{}', siteId })
      const responseMsg = JSON.stringify({ kind: 'allCamera', cameras })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.allCamera(context, siteId).then(result => {
        result.should.eql(cameras)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('getCamera', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-camera')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'getCamera', span: '{}', id })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.getCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'getCamera', span: '{}', id })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'getCamera', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'bbbb-bbbb'
      const camera = { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      const requestMsg = JSON.stringify({ kind: 'getCamera', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'getCamera', camera })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.getCamera(context, id).then(result => {
        result.should.eql(camera)
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('updateCamera', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('update-camera')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'bbbb-bbbb'
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
      const requestMsg = JSON.stringify({ kind: 'updateCamera', span: '{}', id, input })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.updateCamera(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'bbbb-bbbb'
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
      const requestMsg = JSON.stringify({ kind: 'updateCamera', span: '{}', id, input })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateCamera(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'bbbb-bbbb'
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
      const requestMsg = JSON.stringify({ kind: 'updateCamera', span: '{}', id, input })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateCamera(context, id, input).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'bbbb-bbbb'
      const input = { siteId: '1111-1111', serialNo: '2001', resolution: 1920000 }
      const requestMsg = JSON.stringify({ kind: 'updateCamera', span: '{}', id, input })
      const responseMsg = JSON.stringify({ kind: 'updateCamera', updated: true })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.updateCamera(context, id, input).then(result => {
        result.should.be.true()
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })

  describe('deleteCamera', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('delete-camera')

        const spanContext = `{}`
        return func(spanContext)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with error when request is timed out', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'deleteCamera', span: '{}', id })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(new nats.NatsError('timeout'))
      })
      service.deleteCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response is not json', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'deleteCamera', span: '{}', id })
      const responseMsg = `{ invalid json }`
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should reject with error when response has an expected type', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'deleteCamera', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'invalidKind' })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteCamera(context, id).catch(err => {
        should.exist(err)
        _conn.verify()
        _conn.restore()
        done()
      })
    })

    it('should resolve successfully when response has the expected type', done => {
      const id = 'bbbb-bbbb'
      const requestMsg = JSON.stringify({ kind: 'deleteCamera', span: '{}', id })
      const responseMsg = JSON.stringify({ kind: 'deleteCamera', deleted: true })
      _conn.expects('requestOne').callsFake((subject, message, options, timeout, callback) => {
        subject.should.equal(expectedSubject)
        message.should.equal(requestMsg)
        options.should.eql(expectedOptions)
        timeout.should.equal(expectedTimeout)
        callback(responseMsg)
      })
      service.deleteCamera(context, id).then(result => {
        result.should.be.true()
        _conn.verify()
        _conn.restore()
        done()
      }).catch(done)
    })
  })
})
