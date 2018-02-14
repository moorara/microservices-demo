/* eslint-env mocha */
require('should-http')
const should = require('should')
const rp = require('request-promise')

const serviceUrl = process.env.SERVICE_URL || 'http://localhost:4010'

describe('site-service', () => {
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
        res.body.should.match(/# TYPE nodejs_active_handles_total gauge/)
        res.body.should.match(/# TYPE nodejs_active_requests_total gauge/)
        res.body.should.match(/# TYPE nodejs_heap_size_total_bytes gauge/)
        res.body.should.match(/# TYPE nodejs_heap_size_used_bytes gauge/)
        res.body.should.match(/# TYPE process_max_fds gauge/)
        res.body.should.match(/# TYPE process_open_fds gauge/)
        res.body.should.match(/# TYPE process_heap_bytes gauge/)
        res.body.should.match(/# TYPE process_resident_memory_bytes gauge/)
        res.body.should.match(/# TYPE process_virtual_memory_bytes gauge/)
        res.body.should.match(/# TYPE http_requests_duration_seconds histogram/)
        res.body.should.match(/# TYPE http_requests_duration_quantiles_seconds summary/)
        done()
      }).catch(done)
    })
  })

  describe('site with name and location', () => {
    let id
    const name = 'New Site'
    const location = 'Ottawa, ON'
    const newName = 'Plant Site'
    const newLocation = 'Ottawa, ON, CANADA'

    it('POST /v1/sites', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/sites`
      opts.body = { name, location }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        done()
      }).catch(done)
    })

    it('GET /v1/sites', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].name.should.equal(name)
        res.body[0].location.should.equal(location)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        done()
      }).catch(done)
    })

    it('PATCH /v1/sites/:id', done => {
      opts.method = 'PATCH'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(location)
        done()
      }).catch(done)
    })

    it('PUT /v1/sites/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName, location: newLocation }
      rp(opts).then(res => {
        res.should.have.status(204)
        should.not.exist(res.body)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        done()
      }).catch(done)
    })

    it('DELETE /v1/sites/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        done()
      }).catch(done)
    })
  })

  describe('site with name, location, and priority', () => {
    let id
    const name = 'New Site'
    const location = 'Ottawa, ON'
    const priority = 3
    const newName = 'Plant Site'
    const newLocation = 'Ottawa, ON, CANADA'
    const newPriority = 2

    it('POST /v1/sites', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/sites`
      opts.body = { name, location, priority }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('GET /v1/sites', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].name.should.equal(name)
        res.body[0].location.should.equal(location)
        res.body[0].priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('PATCH /v1/sites/:id', done => {
      opts.method = 'PATCH'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName, location: newLocation }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('PUT /v1/sites/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName, location: newLocation, priority: newPriority }
      rp(opts).then(res => {
        res.should.have.status(204)
        should.not.exist(res.body)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.priority.should.equal(newPriority)
        done()
      }).catch(done)
    })

    it('DELETE /v1/sites/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.priority.should.equal(newPriority)
        done()
      }).catch(done)
    })
  })

  describe('site with name, location, tags, and priority', () => {
    let id
    const name = 'New Site'
    const location = 'Ottawa, ON'
    const tags = ['hydro', 'power']
    const priority = 3
    const newName = 'Plant Site'
    const newLocation = 'Ottawa, ON, CANADA'
    const newTags = ['hydro', 'power', 'plant']
    const newPriority = 2

    it('POST /v1/sites', done => {
      opts.method = 'POST'
      opts.uri = `${serviceUrl}/v1/sites`
      opts.body = { name, location, tags, priority }
      rp(opts).then(res => {
        id = res.body.id
        res.should.have.status(201)
        should.exist(res.body.id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        res.body.tags.should.eql(tags)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('GET /v1/sites', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.should.have.length(1)
        res.body[0].id.should.equal(id)
        res.body[0].name.should.equal(name)
        res.body[0].location.should.equal(location)
        res.body[0].tags.should.eql(tags)
        res.body[0].priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(name)
        res.body.location.should.equal(location)
        res.body.tags.should.eql(tags)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('PATCH /v1/sites/:id', done => {
      opts.method = 'PATCH'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName, location: newLocation, tags: newTags }
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.tags.should.eql(newTags)
        res.body.priority.should.equal(priority)
        done()
      }).catch(done)
    })

    it('PUT /v1/sites/:id', done => {
      opts.method = 'PUT'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      opts.body = { name: newName, location: newLocation, tags: newTags, priority: newPriority }
      rp(opts).then(res => {
        res.should.have.status(204)
        should.not.exist(res.body)
        done()
      }).catch(done)
    })

    it('GET /v1/sites/:id', done => {
      opts.method = 'GET'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.tags.should.eql(newTags)
        res.body.priority.should.equal(newPriority)
        done()
      }).catch(done)
    })

    it('DELETE /v1/sites/:id', done => {
      opts.method = 'DELETE'
      opts.uri = `${serviceUrl}/v1/sites/${id}`
      rp(opts).then(res => {
        res.should.have.status(200)
        res.body.id.should.equal(id)
        res.body.name.should.equal(newName)
        res.body.location.should.equal(newLocation)
        res.body.tags.should.eql(newTags)
        res.body.priority.should.equal(newPriority)
        done()
      }).catch(done)
    })
  })
})
