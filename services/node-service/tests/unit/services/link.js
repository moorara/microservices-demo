/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')

const LinkService = require('../../../services/link')

describe('LinkService', () => {
  let config, logger
  let Model, _Model
  let modelInstance, _modelInstance
  let service, _service

  beforeEach(() => {
    config = {}
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    Model = function (data) {
      return modelInstance
    }
    Model.limit = () => {}
    Model.skip = () => {}
    Model.exec = () => {}
    Model.find = () => {}
    Model.findById = () => {}
    Model.findByIdAndUpdate = () => {}
    Model.findByIdAndRemove = () => {}
    _Model = sinon.mock(Model)

    modelInstance = {}
    modelInstance.save = () => {}
    _modelInstance = sinon.mock(modelInstance)

    service = new LinkService(config, {
      logger,
      LinkModel: Model
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
      service = new LinkService({})
      should.exist(service.logger)
      should.exist(service.LinkModel)
    })
  })

  describe('create', () => {
    let specs
    let savedLink

    beforeEach(() => {
      specs = {
        url: 'https://nodejs.org',
        title: 'Node.js',
        tags: ['JavaScript'],
        rank: 1
      }
      savedLink = Object.assign({}, specs, {
        id: '2222-bbbb-4444-dddd',
        createdAt: new Date(),
        updatedAt: new Date()
      })
      savedLink.toJSON = () => savedLink
    })

    it('should reject with error when model query fails', done => {
      _modelInstance.expects('save').rejects(new Error('error'))
      service.create(specs).catch(err => {
        _modelInstance.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolve with new link when model query succeeds', done => {
      _modelInstance.expects('save').resolves(savedLink)
      service.create(specs).then(t => {
        _modelInstance.verify()
        t.should.eql(savedLink)
        done()
      }).catch(done)
    })
  })

  describe('getAll', () => {
    let query
    let t1, t2, t3

    beforeEach(() => {
      query = {}

      t1 = { id: '1111-aaaa' }
      t1.toJSON = () => t1

      t2 = { id: '2222-bbbb' }
      t2.toJSON = () => t2

      t3 = { id: '3333-cccc' }
      t3.toJSON = () => t3
    })

    it('should rejects with error when model query fails', done => {
      _Model.expects('find').returns(Model)
      _Model.expects('limit').returns(Model)
      _Model.expects('skip').returns(Model)
      _Model.expects('exec').rejects(new Error('error'))
      service.getAll().catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolves with links when model query succeeds', done => {
      _Model.expects('find').returns(Model)
      _Model.expects('limit').returns(Model)
      _Model.expects('skip').returns(Model)
      _Model.expects('exec').resolves([ t1, t2, t3 ])
      service.getAll().then(links => {
        _Model.verify()
        links.should.eql([ t1, t2, t3 ])
        done()
      }).catch(done)
    })
    it('should resolves with links when model query succeeds', done => {
      query = { url: 'com', title: 'website', tags: 'javascript,go', minRank: '3', maxRank: '10', limit: '20', skip: '10' }
      let mongoQuery = { url: /.*com.*/i, title: /.*website.*/i, tags: { $in: ['javascript', 'go'] }, rank: { $gte: 3, $lte: 10 } }
      _Model.expects('find').withArgs(mongoQuery).returns(Model)
      _Model.expects('limit').withArgs(+query.limit).returns(Model)
      _Model.expects('skip').withArgs(+query.skip).returns(Model)
      _Model.expects('exec').resolves([ t1, t2, t3 ])
      service.getAll(query).then(links => {
        _Model.verify()
        links.should.eql([ t1, t2, t3 ])
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    let id, link

    beforeEach(() => {
      id = '1111-aaaa'
      link = { id }
      link.toJSON = () => link
    })

    it('should rejects with error when model query fails', done => {
      _Model.expects('findById').withArgs(id).rejects(new Error('error'))
      service.get(id).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolves with empty result when model query returns no result', done => {
      _Model.expects('findById').withArgs(id).resolves()
      service.get(id).then(t => {
        _Model.verify()
        should.not.exist(t)
        done()
      }).catch(done)
    })
    it('should resolves with link when model query succeeds', done => {
      _Model.expects('findById').withArgs(id).resolves(link)
      service.get(id).then(t => {
        _Model.verify()
        t.id.should.equal(id)
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    let id, specs, link

    beforeEach(() => {
      id = '2222-bbbb'
      specs = { url: 'https://nodejs.org', title: 'Node.js', tags: ['javascript'], rank: 1 }
      link = Object.assign({ id }, specs)
      link.toJSON = () => link
    })

    it('should rejects with error when model query fails', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs).rejects(new Error('error'))
      service.update(id, specs).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolves with empty result when model query returns no result', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs).resolves()
      service.update(id, specs).then(t => {
        _Model.verify()
        should.not.exist(t)
        done()
      }).catch(done)
    })
    it('should resolves with updated link when model query succeeds', done => {
      _Model.expects('findByIdAndUpdate').withArgs(id, specs).resolves(link)
      service.update(id, specs).then(t => {
        _Model.verify()
        t.id.should.equal(id)
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    let id, link

    beforeEach(() => {
      id = '3333-cccc'
      link = { id }
      link.toJSON = () => link
    })

    it('should rejects with error when model query fails', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).rejects(new Error('error'))
      service.delete(id).catch(err => {
        _Model.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should resolves with empty result when model query returns no result', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).resolves()
      service.delete(id).then(t => {
        _Model.verify()
        should.not.exist(t)
        done()
      }).catch(done)
    })
    it('should resolves with deleted link when model query succeeds', done => {
      _Model.expects('findByIdAndRemove').withArgs(id).resolves(link)
      service.delete(id).then(t => {
        _Model.verify()
        t.id.should.equal(id)
        done()
      }).catch(done)
    })
  })
})
