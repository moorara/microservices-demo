/* eslint-env mocha */
const express = require('express')
const sinon = require('sinon')
const should = require('should')
const supertest = require('supertest')

const LinkRouter = require('../../../routes/link')

describe('LinkRouter', () => {
  let config, logger
  let linkService, _linkService
  let router, _router
  let app, request

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

    linkService = {
      create () {},
      getAll () {},
      get () {},
      update () {},
      delete () {}
    }
    _linkService = sinon.mock(linkService)

    router = new LinkRouter(config, {
      logger,
      linkService
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
    _linkService.restore()
    _router.restore()
  })

  describe('constructor', () => {
    it('should create a new router with defaults', () => {
      router = new LinkRouter({})
      should.exist(router.router)
      should.exist(router.logger)
      should.exist(router.linkService)
    })
  })

  describe('POST /links', () => {
    let specs, link

    beforeEach(() => {
      specs = { url: 'https://nodejs.org', title: 'Node.js', tags: ['javascript'], rank: 1 }
      link = Object.assign({ id: '1111-aaaa' }, specs)
    })

    it('201 - create a link', done => {
      _linkService.expects('create').withArgs(specs).resolves(link)
      request.post('/')
        .send(specs)
        .expect(201, link, done)
    })
    it('415 - error when body is not json', done => {
      request.post('/')
        .send('invalid json')
        .expect(415, done)
    })
    it('500 - error when service fails', done => {
      _linkService.expects('create').withArgs(specs).rejects(new Error('error'))
      request.post('/')
        .send(specs)
        .expect(500, done)
    })
  })

  describe('GET /links', () => {
    let query, links

    beforeEach(() => {
      query = { url: 'com', title: 'website', tags: 'javascript,go', minRank: '1', maxRank: '10', limit: '20', skip: '10' }
      links = [
        { id: '1111-aaaa', url: 'https://golang.org', title: 'Golang Website', tags: ['go'], rank: 3 },
        { id: '2222-bbbb', url: 'https://github.com', title: 'GitHub Website', tags: ['git', 'go'], rank: 5 },
        { id: '3333-cccc', url: 'https://www.docker.com', title: 'Docker Website', tags: ['docker'], rank: 7 }
      ]
    })

    it('200 - get links', done => {
      _linkService.expects('getAll').withArgs(query).resolves(links)
      request.get('/?url=com&title=website&tags=javascript,go&minRank=1&maxRank=10&limit=20&skip=10')
        .expect(200, links, done)
    })
    it('500 - error when service fails', done => {
      _linkService.expects('getAll').withArgs(query).rejects(new Error('error'))
      request.get('/?url=com&title=website&tags=javascript,go&minRank=1&maxRank=10&limit=20&skip=10')
        .expect(500, done)
    })
  })

  describe('GET /links/:id', () => {
    let id, link

    beforeEach(() => {
      id = '1111-aaaa'
      link = { id, url: 'https://nodejs.org', title: 'Node.js', tags: ['javascript'], rank: 1 }
    })

    it('200 - get a link', done => {
      _linkService.expects('get').withArgs(id).resolves(link)
      request.get(`/${id}`)
        .expect(200, link, done)
    })
    it('404 - link not found', done => {
      _linkService.expects('get').withArgs(id).resolves()
      request.get(`/${id}`)
        .expect(404, done)
    })
    it('500 - error when service fails', done => {
      _linkService.expects('get').withArgs(id).rejects(new Error('error'))
      request.get(`/${id}`)
        .expect(500, done)
    })
  })

  describe('PUT /links/:id', () => {
    let id, specs, link

    beforeEach(() => {
      id = '2222-bbbb'
      specs = { url: 'https://nodejs.org', title: 'Node.js', tags: ['javascript'], rank: 1 }
      link = Object.assign({ id }, specs)
    })

    it('200 - update a link', done => {
      _linkService.expects('update').withArgs(id, specs).resolves(link)
      request.put(`/${id}`)
        .send(specs)
        .expect(200, link, done)
    })
    it('404 - link not found', done => {
      _linkService.expects('update').withArgs(id, specs).resolves()
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
      _linkService.expects('update').withArgs(id, specs).rejects(new Error('error'))
      request.put(`/${id}`)
        .send(specs)
        .expect(500, done)
    })
  })

  describe('DELETE /links/:id', () => {
    let id, link

    beforeEach(() => {
      id = '3333-cccc'
      link = { id, url: 'https://nodejs.org', title: 'Node.js', tags: ['javascript'], rank: 1 }
    })

    it('200 - delete a link', done => {
      _linkService.expects('delete').withArgs(id).resolves(link)
      request.delete(`/${id}`)
        .expect(200, link, done)
    })
    it('404 - link not found', done => {
      _linkService.expects('delete').withArgs(id).resolves()
      request.delete(`/${id}`)
        .expect(404, done)
    })
    it('500 - error when service fails', done => {
      _linkService.expects('delete').withArgs(id).rejects(new Error('error'))
      request.delete(`/${id}`)
        .expect(500, done)
    })
  })
})
