/* eslint-env mocha */
const express = require('express')
const sinon = require('sinon')
const should = require('should')
const supertest = require('supertest')

const SiteRouter = require('../../../routes/site')

describe('SiteRouter', () => {
  let config, logger
  let siteService, _siteService
  let router, _router
  let app, request

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

    siteService = {
      create () {},
      all () {},
      get () {},
      update () {},
      modify () {},
      delete () {}
    }
    _siteService = sinon.mock(siteService)

    router = new SiteRouter(config, {
      logger,
      siteService
    })
    _router = sinon.mock(router)

    app = express()
    app.use('/', router.router)
    app.use((err, req, res, next) => {
      // Make express stop logging errors
      if (!res.headersSent) {
        res.status(500).send(err)
      }
    })

    request = supertest(app)
  })

  afterEach(() => {
    _siteService.restore()
    _router.restore()
  })

  describe('constructor', () => {
    it('should create a new router with defaults', () => {
      router = new SiteRouter({})
      should.exist(router.router)
      should.exist(router.logger)
      should.exist(router.siteService)
    })
  })

  describe('POST /sites', () => {
    let specs, site

    beforeEach(() => {
      specs = { name: 'New Site', location: 'Ottawa, ON', tags: ['hydro', 'power'], priority: 3 }
      site = Object.assign({ id: '1111-aaaa' }, specs)
    })

    it('201 - create a site', done => {
      _siteService.expects('create').withArgs(specs).resolves(site)
      request.post('/')
        .send(specs)
        .expect(201, site, done)
    })
    it('415 - error when body is not json', done => {
      request.post('/')
        .send('invalid json')
        .expect(415, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('create').withArgs(specs).rejects(new Error('error'))
      request.post('/')
        .send(specs)
        .expect(500, done)
    })
  })

  describe('GET /sites', () => {
    let query, sites

    beforeEach(() => {
      query = { name: 'Site', location: 'Ottawa', tags: 'hydro,power', minPriority: '2', maxPriority: '4', limit: '10', skip: '10' }
      sites = [
        { id: '1111-aaaa', name: 'Old Site', location: 'Waterloo, ON', tags: ['hydro'], priority: 2 },
        { id: '2222-bbbb', name: 'New Site', location: 'Ottawa, ON', tags: ['hydro', 'power'], priority: 3 },
        { id: '3333-cccc', name: 'Future Site', location: 'Vancouver, BC', tags: ['power'], priority: 4 }
      ]
    })

    it('200 - get sites', done => {
      _siteService.expects('all').withArgs(query).resolves(sites)
      request.get('/?name=Site&location=Ottawa&tags=hydro,power&minPriority=2&maxPriority=4&limit=10&skip=10')
        .expect(200, sites, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('all').withArgs(query).rejects(new Error('error'))
      request.get('/?name=Site&location=Ottawa&tags=hydro,power&minPriority=2&maxPriority=4&limit=10&skip=10')
        .expect(500, done)
    })
  })

  describe('GET /sites/:id', () => {
    let id, site

    beforeEach(() => {
      id = '1111-aaaa'
      site = { id, name: 'New Site', location: 'Ottawa, ON', tags: ['hydro', 'power'], priority: 3 }
    })

    it('200 - get a site', done => {
      _siteService.expects('get').withArgs(id).resolves(site)
      request.get(`/${id}`)
        .expect(200, site, done)
    })
    it('404 - site not found', done => {
      _siteService.expects('get').withArgs(id).resolves()
      request.get(`/${id}`)
        .expect(404, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('get').withArgs(id).rejects(new Error('error'))
      request.get(`/${id}`)
        .expect(500, done)
    })
  })

  describe('PUT /sites/:id', () => {
    let id, specs

    beforeEach(() => {
      id = '1111-aaaa'
      specs = { name: 'Plant Site', location: 'Ottawa, ON, CANADA', priority: 2 }
    })

    it('204 - update a site', done => {
      _siteService.expects('update').withArgs(id, specs).resolves(true)
      request.put(`/${id}`)
        .send(specs)
        .expect(204, done)
    })
    it('404 - site not found', done => {
      _siteService.expects('update').withArgs(id, specs).resolves(false)
      request.put(`/${id}`)
        .send(specs)
        .expect(404, done)
    })
    it('415 - error when body is not json', done => {
      request.put(`/${id}`)
        .send('invalid json')
        .expect(415, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('update').withArgs(id, specs).rejects(new Error('error'))
      request.put(`/${id}`)
        .send(specs)
        .expect(500, done)
    })
  })

  describe('PATCH /sites/:id', () => {
    let id, specs
    let site

    beforeEach(() => {
      id = '1111-aaaa'
      specs = { name: 'Plant Site', location: 'Ottawa, ON, CANADA' }
      site = Object.assign({ id }, specs)
    })

    it('200 - modify a site', done => {
      _siteService.expects('modify').withArgs(id, specs).resolves(site)
      request.patch(`/${id}`)
        .send(specs)
        .expect(200, site, done)
    })
    it('404 - site not found', done => {
      _siteService.expects('modify').withArgs(id, specs).resolves()
      request.patch(`/${id}`)
        .send(specs)
        .expect(404, done)
    })
    it('415 - error when body is not json', done => {
      request.patch(`/${id}`)
        .send('invalid json')
        .expect(415, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('modify').withArgs(id, specs).rejects(new Error('error'))
      request.patch(`/${id}`)
        .send(specs)
        .expect(500, done)
    })
  })

  describe('DELETE /sites/:id', () => {
    let id, site

    beforeEach(() => {
      id = '1111-aaaa'
      site = { id, name: 'New Site', location: 'Ottawa, ON', tags: ['hydro', 'power'], priority: 3 }
    })

    it('200 - delete a site', done => {
      _siteService.expects('delete').withArgs(id).resolves(site)
      request.delete(`/${id}`)
        .expect(200, site, done)
    })
    it('404 - site not found', done => {
      _siteService.expects('delete').withArgs(id).resolves()
      request.delete(`/${id}`)
        .expect(404, done)
    })
    it('500 - error when service fails', done => {
      _siteService.expects('delete').withArgs(id).rejects(new Error('error'))
      request.delete(`/${id}`)
        .expect(500, done)
    })
  })
})
