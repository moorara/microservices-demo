/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const SiteService = require('./site')

describe('SiteService', () => {
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
        baseUrl: 'http://localhost/'
      }
    }
    _axios = sinon.mock(axios)

    config = {}
    options = { logger, histogram, summary, tracer, axios }

    service = new SiteService(config, options)
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
      const service = new SiteService(config, { tracer: options.tracer })
      should.exist(service.logger)
      should.exist(service.histogram)
      should.exist(service.summary)
      should.exist(service.tracer)
      should.exist(service.axios)
    })
    it('should create a new service with provided options', () => {
      const service = new SiteService(config, options)
      service.logger.should.equal(options.logger)
      service.histogram.should.equal(options.histogram)
      service.summary.should.equal(options.summary)
      service.tracer.should.equal(options.tracer)
      service.axios.should.equal(options.axios)
    })
  })

  describe('exec', () => {
    const verifyTrace = (spanName, logMessage) => {
      const span = tracer._spans[0]
      span._operationName.should.equal(spanName)
      span._tags[opentracing.Tags.SPAN_KIND].should.equal('client')
      span._tags[opentracing.Tags.PEER_SERVICE].should.equal('site-service')
      span._tags[opentracing.Tags.PEER_ADDRESS].should.equal(axios.defaults.baseUrl)
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
        name.should.equal('create-site')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const input = { name: 'Plant', location: 'Canada', priority: 1, tags: ['energy'] }
      const err = new Error('create error')
      _axios.expects('request').withArgs({ headers: {}, method: 'post', url: '/sites', data: input }).rejects(err)
      service.create(context, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const input = { name: 'Plant', location: 'Canada', priority: 1, tags: ['energy'] }
      const site = Object.assign({}, { id: 'aaaa-aaaa' }, input)
      _axios.expects('request').withArgs({ headers: {}, method: 'post', url: '/sites', data: input }).resolves({ data: site })
      service.create(context, input).then(s => {
        s.should.equal(site)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-sites')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const query = { name: 'Plant', location: 'Canada', tags: 'energy', minPriority: 2, maxPriority: 4, limit: 10, skip: 10 }
      const err = new Error('all error')
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: '/sites', params: query }).rejects(err)
      service.all(context, query).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const query = { name: 'Plant', location: 'Canada', tags: 'energy', minPriority: 2, maxPriority: 4, limit: 10, skip: 10 }
      const sites = [
        { name: 'Power Plant', location: 'Toronto, Canada', tags: [ 'energy', 'power' ], priority: 2 },
        { name: 'Hydro Plant', location: 'Montreal, Canada', tags: [ 'energy', 'hydro' ], priority: 4 }
      ]
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: '/sites', params: query }).resolves({ data: sites })
      service.all(context, query).then(s => {
        s.should.equal(sites)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('get-site')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const id = 'aaaa-aaaa'
      const err = new Error('get error')
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sites/${id}` }).rejects(err)
      service.get(context, id).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = 'aaaa-aaaa'
      const site = { id: 'aaaa-aaaa', name: 'Plant', location: 'Canada', priority: 1, tags: ['energy'] }
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sites/${id}` }).resolves({ data: site })
      service.get(context, id).then(s => {
        s.should.equal(site)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('update-site')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when first request rejects', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Power Plant', location: 'Toronto, Canada', priority: 2, tags: ['energy', 'power'] }
      const err = new Error('update error')
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sites/${id}`, data: input }).rejects(err)
      service.update(context, id, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should reject with an error when second request rejects', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Power Plant', location: 'Toronto, Canada', priority: 2, tags: ['energy', 'power'] }
      const err = new Error('get error')
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sites/${id}`, data: input }).resolves({})
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sites/${id}` }).rejects(err)
      service.update(context, id, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Power Plant', location: 'Toronto, Canada', priority: 2, tags: ['energy', 'power'] }
      const site = { id: 'aaaa-aaaa', name: 'Power Plant', location: 'Toronto, Canada', priority: 2, tags: ['energy', 'power'] }
      _axios.expects('request').withArgs({ headers: {}, method: 'put', url: `/sites/${id}`, data: input }).resolves({})
      _axios.expects('request').withArgs({ headers: {}, method: 'get', url: `/sites/${id}` }).resolves({ data: site })
      service.update(context, id, input).then(s => {
        s.should.equal(site)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('modify', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('modify-site')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Power Plant', location: 'Toronto, Canada' }
      const err = new Error('modify error')
      _axios.expects('request').withArgs({ headers: {}, method: 'patch', url: `/sites/${id}`, data: input }).rejects(err)
      service.modify(context, id, input).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Power Plant', location: 'Toronto, Canada' }
      const site = { id: 'aaaa-aaaa', name: 'Power Plant', location: 'Toronto, Canada', priority: 1, tags: ['energy'] }
      _axios.expects('request').withArgs({ headers: {}, method: 'patch', url: `/sites/${id}`, data: input }).resolves({ data: site })
      service.modify(context, id, input).then(s => {
        s.should.equal(site)
        _axios.verify()
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    beforeEach(() => {
      _service.expects('exec').callsFake((ctx, name, func) => {
        ctx.should.eql(context)
        name.should.equal('delete-site')

        const headers = {}
        return func(headers)
      })
    })

    afterEach(() => {
      _service.verify()
      _service.restore()
    })

    it('should reject with an error when request rejects', done => {
      const id = 'aaaa-aaaa'
      const err = new Error('delete error')
      _axios.expects('request').withArgs({ headers: {}, method: 'delete', url: `/sites/${id}` }).rejects(err)
      service.delete(context, id).catch(e => {
        e.should.equal(err)
        _axios.verify()
        done()
      })
    })
    it('should resolve successfully when request resolves', done => {
      const id = 'aaaa-aaaa'
      const result = { data: {} }
      _axios.expects('request').withArgs({ headers: {}, method: 'delete', url: `/sites/${id}` }).resolves(result)
      service.delete(context, id).then(r => {
        r.should.eql(result.data)
        _axios.verify()
        done()
      }).catch(done)
    })
  })
})
