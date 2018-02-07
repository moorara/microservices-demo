/* eslint-env mocha */
require('should-http')
const should = require('should')
const rp = require('request-promise')

const serviceUrl = process.env.SERVICE_URL || 'http://localhost:4020'

describe('node-service', () => {
  let opts

  beforeEach(() => {
    opts = {
      json: true,
      simple: false,
      resolveWithFullResponse: true
    }
  })

  describe('general endpoints', () => {
    it('GET /health', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/health`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.equal('OK')
        done()
      }).catch(done)
    })

    it('GET /metrics', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/metrics`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.match(/# TYPE nodejs_version_info gauge/)
        res.body.should.match(/# TYPE http_requests_duration_seconds histogram/)
        res.body.should.match(/# TYPE http_requests_duration_quantiles_seconds summary/)
        done()
      }).catch(done)
    })
  })

  describe('link with url and title', () => {
    let id
    const url = 'https://nodejs.org'
    const title = 'Node.js'
    const newUrl = 'https://expressjs.com'
    const newTitle = 'Express.js'

    it('POST /v1/links', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/links`
      opts.body = { url, title }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        done()
      }).catch(done)
    })

    it('GET /v1/links', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].url.should.equal(url)
        res.body[0].title.should.equal(title)
        done()
      }).catch(done)
    })

    it('GET /v1/links/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        done()
      }).catch(done)
    })

    it('PUT /v1/links/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      opts.body = { url: newUrl, title: newTitle }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        done()
      }).catch(done)
    })

    it('DELETE /v1/links/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        done()
      }).catch(done)
    })
  })

  describe('link with url, title, and rank', () => {
    let id
    const url = 'https://nodejs.org'
    const title = 'Node.js'
    const rank = 1
    const newUrl = 'https://expressjs.com'
    const newTitle = 'Express.js'
    const newRank = 2

    it('POST /v1/links', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/links`
      opts.body = { url, title, rank }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        res.body.rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('GET /v1/links', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].url.should.equal(url)
        res.body[0].title.should.equal(title)
        res.body[0].rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('GET /v1/links/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        res.body.rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('PUT /v1/links/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      opts.body = { url: newUrl, title: newTitle, rank: newRank }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        res.body.rank.should.equal(newRank)
        done()
      }).catch(done)
    })

    it('DELETE /v1/links/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        res.body.rank.should.equal(newRank)
        done()
      }).catch(done)
    })
  })

  describe('link with url, title, tags, and rank', () => {
    let id
    const url = 'https://nodejs.org'
    const title = 'Node.js'
    const tags = ['javascipt', 'node']
    const rank = 1
    const newUrl = 'https://expressjs.com'
    const newTitle = 'Express.js'
    const newTags = ['javascipt', 'express']
    const newRank = 2

    it('POST /v1/links', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/links`
      opts.body = { url, title, tags, rank }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        res.body.tags.should.eql(tags)
        res.body.rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('GET /v1/links', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].url.should.equal(url)
        res.body[0].title.should.equal(title)
        res.body[0].tags.should.eql(tags)
        res.body[0].rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('GET /v1/links/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(url)
        res.body.title.should.equal(title)
        res.body.tags.should.eql(tags)
        res.body.rank.should.equal(rank)
        done()
      }).catch(done)
    })

    it('PUT /v1/links/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      opts.body = { url: newUrl, title: newTitle, tags: newTags, rank: newRank }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        res.body.tags.should.eql(newTags)
        res.body.rank.should.equal(newRank)
        done()
      }).catch(done)
    })

    it('DELETE /v1/links/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/links/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.url.should.equal(newUrl)
        res.body.title.should.equal(newTitle)
        res.body.tags.should.eql(newTags)
        res.body.rank.should.equal(newRank)
        done()
      }).catch(done)
    })
  })
})
