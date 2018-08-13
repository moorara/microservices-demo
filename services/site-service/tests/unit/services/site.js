/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')
const opentracing = require('opentracing')

const SiteService = require('../../../services/site')

describe('SiteService', () => {
  let config, logger, tracer
  let Model, _Model
  let modelInstance, _modelInstance
  let service, _service

  beforeEach(() => {
    config = {}
    logger = {
      trace () {},
      debug () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }
    tracer = new opentracing.MockTracer()

    Model = function (data) {
      return modelInstance
    }
    Model.limit = () => {}
    Model.skip = () => {}
    Model.exec = () => {}
    Model.find = () => {}
    Model.findById = () => {}
    Model.update = () => {}
    Model.findByIdAndUpdate = () => {}
    Model.findByIdAndRemove = () => {}
    _Model = sinon.mock(Model)

    modelInstance = {}
    modelInstance.save = () => {}
    _modelInstance = sinon.mock(modelInstance)

    service = new SiteService(config, {
      logger,
      tracer,
      SiteModel: Model
    })
    _service = sinon.mock(service)
  })

  afterEach(() => {
    _Model.restore()
    _modelInstance.restore()
    _service.restore()
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      service = new SiteService({})
      should.exist(service.logger)
      should.exist(service.SiteModel)
    })
  })

  describe('create', () => {
    let specs, context
    let savedSite

    beforeEach(() => {
      specs = {
        name: 'New Site',
        location: 'Ottawa, ON',
        tags: ['hydro', 'power'],
        priority: 3
      }
      context = {
        span: tracer.startSpan('mock-span')
      }

      savedSite = Object.assign({}, specs, {
        id: '1111-aaaa',
        createdAt: new Date(),
        updatedAt: new Date()
      })
      savedSite.toJSON = () => savedSite
    })

    it('should reject with error when model query fails', done => {
      _modelInstance.expects('save').rejects(new Error('error'))
      service.create(specs, context).catch(err => {
        _modelInstance.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with new site when model query succeeds', done => {
      _modelInstance.expects('save').resolves(savedSite)
      service.create(specs, context).then(t => {
        _modelInstance.verify()
        t.should.eql(savedSite)
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    let query, context
    let s1, s2, s3

    beforeEach(() => {
      query = {}
      context = { span: tracer.startSpan('mock-span') }

      s1 = { id: '1111-aaaa' }
      s1.toJSON = () => s1

      s2 = { id: '2222-bbbb' }
      s2.toJSON = () => s2

      s3 = { id: '3333-cccc' }
      s3.toJSON = () => s3
    })

    it('should reject with error when model query fails', done => {
      _Model.expects('find').returns(Model)
      _Model.expects('limit').returns(Model)
      _Model.expects('skip').returns(Model)
      _Model.expects('exec').rejects(new Error('error'))
      service.all(query, context).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with sites when model query succeeds', done => {
      _Model.expects('find').returns(Model)
      _Model.expects('limit').returns(Model)
      _Model.expects('skip').returns(Model)
      _Model.expects('exec').resolves([ s1, s2, s3 ])
      service.all(query, context).then(sites => {
        _Model.verify()
        sites.should.eql([ s1, s2, s3 ])
        done()
      }).catch(done)
    })
    it('should resolve with sites when model query succeeds', done => {
      query = { name: 'Site', location: 'Ottawa', tags: 'hydro,power', minPriority: '2', maxPriority: '4', limit: '10', skip: '10' }
      let mongoQuery = { name: /.*Site.*/i, location: /.*Ottawa.*/i, tags: { $in: ['hydro', 'power'] }, priority: { $gte: 2, $lte: 4 } }
      _Model.expects('find').withArgs(mongoQuery).returns(Model)
      _Model.expects('limit').withArgs(+query.limit).returns(Model)
      _Model.expects('skip').withArgs(+query.skip).returns(Model)
      _Model.expects('exec').resolves([ s1, s2, s3 ])
      service.all(query, context).then(sites => {
        _Model.verify()
        sites.should.eql([ s1, s2, s3 ])
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    let id, context
    let site

    beforeEach(() => {
      id = '1111-aaaa'
      context = { span: tracer.startSpan('mock-span') }

      site = { id }
      site.toJSON = () => site
    })

    it('should reject with error when model query fails', done => {
      _Model.expects('findById').withArgs(id).rejects(new Error('error'))
      service.get(id, context).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with empty result when model query returns no result', done => {
      _Model.expects('findById').withArgs(id).resolves()
      service.get(id, context).then(s => {
        _Model.verify()
        should.not.exist(s)
        done()
      }).catch(done)
    })
    it('should resolve with site when model query succeeds', done => {
      _Model.expects('findById').withArgs(id).resolves(site)
      service.get(id, context).then(s => {
        _Model.verify()
        s.id.should.equal(id)
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    let id, specs, context
    let query, options

    beforeEach(() => {
      id = '1111-aaaa'
      specs = { name: 'Plant Site', location: 'Ottawa, ON, CANADA', priority: 2 }
      context = { span: tracer.startSpan('mock-span') }

      query = { _id: id }
      options = { upsert: false, runValidators: true, overwrite: true }
    })

    it('should reject with error when model query fails', done => {
      _Model.expects('update').withArgs(query, specs, options).rejects(new Error('error'))
      service.update(id, specs, context).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with false when model query returns no response', done => {
      _Model.expects('update').withArgs(query, specs, options).resolves({ ok: 1, n: 0, nModified: 0 })
      service.update(id, specs, context).then(r => {
        _Model.verify()
        r.should.be.false()
        done()
      }).catch(done)
    })
    it('should resolve with true when model query returns successful response', done => {
      _Model.expects('update').withArgs(query, specs, options).resolves({ ok: 1, n: 1, nModified: 1 })
      service.update(id, specs, context).then(r => {
        _Model.verify()
        r.should.be.true()
        done()
      }).catch(done)
    })
  })

  describe('modify', () => {
    let id, specs, context
    let options, site

    beforeEach(() => {
      id = '1111-aaaa'
      specs = { name: 'Plant Site', location: 'Ottawa, ON, CANADA' }
      context = { span: tracer.startSpan('mock-span') }

      options = { new: true, upsert: false, runValidators: true }
      site = Object.assign({ id }, specs)
      site.toJSON = () => site
    })

    it('should reject with error when model query fails', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs, options).rejects(new Error('error'))
      service.modify(id, specs, context).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with empty result when model query returns no result', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs, options).resolves()
      service.modify(id, specs, context).then(s => {
        _Model.verify()
        should.not.exist(s)
        done()
      }).catch(done)
    })
    it('should resolve with modified site when model query succeeds', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs, options).resolves(site)
      service.modify(id, specs, context).then(s => {
        _Model.verify()
        s.id.should.equal(id)
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    let id, context
    let site

    beforeEach(() => {
      id = '1111-aaaa'
      context = { span: tracer.startSpan('mock-span') }

      site = { id }
      site.toJSON = () => site
    })

    it('should reject with error when model query fails', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).rejects(new Error('error'))
      service.delete(id, context).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with empty result when model query returns no result', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).resolves()
      service.delete(id, context).then(s => {
        _Model.verify()
        should.not.exist(s)
        done()
      }).catch(done)
    })
    it('should resolve with deleted site when model query succeeds', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).resolves(site)
      service.delete(id, context).then(s => {
        _Model.verify()
        s.id.should.equal(id)
        done()
      }).catch(done)
    })
  })
})
