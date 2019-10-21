/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const SensorService = require('./sensor')

describe('SensorService', () => {
  let logger
  let histogram, _histogram
  let summary, _summary
  let tracer, _tracer
  let axios, _axios
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

    axios = {
      request () {},
      defaults: {
        baseUrl: 'http://sensor-service:4020/v1/'
      }
    }
    _axios = sinon.mock(axios)

    config = {
      sensorServiceAddr: 'sensor-service:4020'
    }
    options = { logger, histogram, summary, tracer, axios }

    service = new SensorService(config, options)
    _service = sinon.mock(service)

    span = {}
    context = { span }
  })

  afterEach(() => {
    _histogram.restore()
    _summary.restore()
    _tracer.restore()
    _axios.restore()
    _service.restore()
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new SensorService(config, { tracer: options.tracer })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
      should.exist(service.axios)
    })
    it('should create a new service with provided options', () => {
      const service = new SensorService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
      service.axios.should.equal(options.axios)
    })
  })

  describe('exec', () => {
    const verifyTrace = (spanName, logMessage) => {
      const span = tracer.report().spans[0]
      span.operationName().should.equal(spanName)
      span.tags()[opentracing.Tags.SPAN_KIND].should.equal('client')
      span.tags()[opentracing.Tags.PEER_SERVICE].should.equal('sensor-service')
      span.tags()[opentracing.Tags.PEER_ADDRESS].should.equal(axios.defaults.baseUrl)
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

  describe('create', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('create-sensor')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
      const err = new Error('create error')
      _axios.expects('request').withArgs({ headers: {}, method: 'post', url: '/sensors', data: input }).rejects(err)
      service.create(context, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
      const sensor = Object.assign({}, { id: '1111-1111' }, input)
      _axios.expects('request').withArgs({ headers: {}, method: 'post', url: '/sensors', data: input }).resolves({ data: sensor })
      service.create(context, input).then(s => {
        s.should.equal(sensor)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-sensors')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const siteId = 'aaaa-aaaa'
      const err = new Error('all error')
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: '/sensors', params: { siteId } }).rejects(err)
      service.all(context, siteId).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const siteId = 'aaaa-aaaa'
      const sensors = [
        { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
        { id: '2222-2222', siteId: 'aaaa-aaaa', name: 'pressure', unit: 'atmosphere', minSafe: 0.5, maxSafe: 1.0 }
      ]
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: '/sensors', params: { siteId } }).resolves({ data: sensors })
      service.all(context, siteId).then(s => {
        s.should.equal(sensors)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-sensor')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const id = '1111-1111'
      const err = new Error('get error')
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sensors/${id}` }).rejects(err)
      service.get(context, id).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = '1111-1111'
      const sensor = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sensors/${id}` }).resolves({ data: sensor })
      service.get(context, id).then(s => {
        s.should.equal(sensor)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('update-sensor')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when first request rejects', done => {
      const id = '1111-1111'
      const input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      const err = new Error('update error')
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sensors/${id}`, data: input }).rejects(err)
      service.update(context, id, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should reject with an error when second request rejects', done => {
      const id = '1111-1111'
      const input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      const err = new Error('get error')
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sensors/${id}`, data: input }).resolves({})
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sensors/${id}` }).rejects(err)
      service.update(context, id, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when requests resolve', done => {
      const id = '1111-1111'
      const input = { siteId: 'aaaa-aaaa', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      const sensor = { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sensors/${id}`, data: input }).resolves({})
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sensors/${id}` }).resolves({ data: sensor })
      service.update(context, id, input).then(s => {
        s.should.equal(sensor)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('delete-sensor')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const id = '1111-1111'
      const err = new Error('delete error')
      _axios.expects('request').withArgs({ headers: {}, method: 'delete', url: `/sensors/${id}` }).rejects(err)
      service.delete(context, id).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = '1111-1111'
      const result = { data: {} }
      _axios.expects('request').withArgs({ headers: {}, method: 'delete', url: `/sensors/${id}` }).resolves(result)
      service.delete(context, id).then(r => {
        r.should.eql(result.data)
        _axios.verify()
        done()
      }).catch(done)
    })
  })
})
