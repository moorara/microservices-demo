/* eslint-env mocha */
const stream = require('stream')
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const SwitchService = require('./switch')

describe('SwitchService', () => {
  let logger
  let histogram, _histogram
  let summary, _summary
  let tracer, _tracer
  let client, _client
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

    client = {
      installSwitch () {},
      removeSwitch () {},
      getSwitch () {},
      getSwitches () {},
      setSwitch () {}
    }
    _client = sinon.mock(client)

    config = {
      switchServiceAddr: 'switch-service:4030'
    }
    options = { logger, histogram, summary, tracer, client }

    service = new SwitchService(config, options)
    _service = sinon.mock(service)

    span = {}
    context = { span }
  })

  afterEach(() => {
    _histogram.restore()
    _summary.restore()
    _tracer.restore()
    _client.restore()
    _service.restore()
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new SwitchService(config, { tracer: options.tracer })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
      should.exist(service.client)
    })
    it('should create a new service with provided options', () => {
      const service = new SwitchService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
      service.client.should.equal(options.client)
    })
  })

  describe('exec', () => {
    const verifyTrace = (spanName, logMessage) => {
      const span = tracer.report().spans[0]
      span.operationName().should.equal(spanName)
      span.tags()[opentracing.Tags.SPAN_KIND].should.equal('client')
      span.tags()[opentracing.Tags.PEER_SERVICE].should.equal('switch-service')
      span.tags()[opentracing.Tags.PEER_ADDRESS].should.equal(config.switchServiceAddr)
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

  describe('installSwitch', () => {
    const metadata = {}

    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('install-switch')
        return func(metadata)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request fails', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'light', state: 'off', states: ['off', 'on'] }
      const err = new Error('install error')
      _client.expects('installSwitch').withArgs(input, metadata).yieldsAsync(err)
      service.installSwitch(context, input).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })
    })
    it('should resolve successfully when request succeeds', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'light', state: 'off', states: ['off', 'on'] }
      const swtch = Object.assign({}, { id: '1111-1111' }, input)
      _client.expects('installSwitch').withArgs(input, metadata).yieldsAsync(null, swtch)
      service.installSwitch(context, input).then(s => {
        s.should.equal(swtch)
        _client.verify()
        done()
      }).catch(done)
    })
  })

  describe('removeSwitch', () => {
    const metadata = {}

    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('remove-switch')
        return func(metadata)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request fails', done => {
      const id = '1111-1111'
      const err = new Error('remove error')
      _client.expects('removeSwitch').withArgs({ id }, metadata).yieldsAsync(err)
      service.removeSwitch(context, id).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })
    })
    it('should resolve successfully when request succeeds', done => {
      const id = '1111-1111'
      const response = {}
      _client.expects('removeSwitch').withArgs({ id }, metadata).yieldsAsync(null, response)
      service.removeSwitch(context, id).then(() => {
        _client.verify()
        done()
      }).catch(done)
    })
  })

  describe('getSwitch', () => {
    const metadata = {}

    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-switch')
        return func(metadata)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request fails', done => {
      const id = '1111-1111'
      const err = new Error('get error')
      _client.expects('getSwitch').withArgs({ id }, metadata).yieldsAsync(err)
      service.getSwitch(context, id).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })
    })
    it('should resolve successfully when request succeeds', done => {
      const id = '1111-1111'
      const swtch = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'light', state: 'off', states: ['off', 'on'] }
      _client.expects('getSwitch').withArgs({ id }, metadata).yieldsAsync(null, swtch)
      service.getSwitch(context, id).then(s => {
        s.should.equal(swtch)
        _client.verify()
        done()
      }).catch(done)
    })
  })

  describe('getSwitches', () => {
    let mockStream
    const metadata = {}

    beforeEach(() => {
      mockStream = new stream.Readable()
      mockStream._read = () => {}

      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-switches')
        return func(metadata)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when response stream fails', done => {
      const siteId = 'aaaa-aaaa'
      const err = new Error('get error')
      _client.expects('getSwitches').withArgs({ siteId }, metadata).returns(mockStream)
      service.getSwitches(context, siteId).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })

      mockStream.emit('error', err)
    })
    it('should resolve successfully when response stream succeeds', done => {
      const siteId = 'aaaa-aaaa'
      const switches = [
        { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'light', state: 'off', states: ['off', 'on'] },
        { id: '2222-2222', siteId: 'aaaa-aaaa', name: 'light', state: 'off', states: ['off', 'on'] }
      ]
      _client.expects('getSwitches').withArgs({ siteId }, metadata).returns(mockStream)
      service.getSwitches(context, siteId).then(s => {
        s.should.eql(switches)
        _client.verify()
        done()
      }).catch(done)

      switches.forEach(swtch => mockStream.emit('data', swtch))
      mockStream.emit('end')
    })
  })

  describe('setSwitch', () => {
    const metadata = {}

    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('set-switch')
        return func(metadata)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when first request fails', done => {
      const id = '1111-1111'
      const state = 'on'
      const err = new Error('set error')
      _client.expects('setSwitch').withArgs({ id, state }, metadata).yieldsAsync(err)
      service.setSwitch(context, id, { state }).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })
    })
    it('should reject with an error when second request fails', done => {
      const id = '1111-1111'
      const state = 'on'
      const err = new Error('get error')
      _client.expects('setSwitch').withArgs({ id, state }, metadata).yieldsAsync(null, {})
      _client.expects('getSwitch').withArgs({ id }, metadata).yieldsAsync(err)
      service.setSwitch(context, id, { state }).catch(e => {
        e.should.equal(err)
        _client.verify()
        done()
      })
    })
    it('should resolve successfully when requests succeed', done => {
      const id = '1111-1111'
      const state = 'on'
      const swtch = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'light', state: 'on', states: ['off', 'on'] }
      _client.expects('setSwitch').withArgs({ id, state }, metadata).yieldsAsync(null, {})
      _client.expects('getSwitch').withArgs({ id }, metadata).yieldsAsync(null, swtch)
      service.setSwitch(context, id, { state }).then(s => {
        s.should.equal(swtch)
        _client.verify()
        done()
      }).catch(done)
    })
  })
})
