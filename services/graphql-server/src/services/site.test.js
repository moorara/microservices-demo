/* eslint-env mocha */
const should = require('should')

const SiteService = require('./site')

describe('SiteService', () => {
  let config, logger

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
  })

  describe('constructor', () => {
    it('should create a new service with defaults', () => {
      const service = new SiteService(config)
      should.exist(service.logger)
    })
    it('should create a new service with provided options', () => {
      const options = { logger }
      const service = new SiteService(config, options)
      service.logger.should.equal(options.logger)
    })
  })

  describe('create', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should create and persist a new site', done => {
      const input = { name: 'Oil Platform', location: 'Vancouver, BC', tags: ['energy', 'oil'] }
      service.create(context, input).then(site => {
        should.exist(site.id)
        site.name.should.equal(input.name)
        site.location.should.equal(input.location)
        site.tags.should.eql(input.tags)
        done()
      }).catch(done)
    })
    it('should create and persist a new site', done => {
      const input = { name: 'Hydropower Plant', location: 'Ottawa, ON', tags: ['energy', 'power', 'hydro'] }
      service.create(context, input).then(site => {
        should.exist(site.id)
        site.name.should.equal(input.name)
        site.location.should.equal(input.location)
        site.tags.should.eql(input.tags)
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should return all sites', done => {
      service.all(context).then(sites => {
        sites.should.have.length(2)
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should return a site by id', done => {
      const id = 'aaaa-aaaa'
      service.get(context, id).then(site => {
        site.id.should.equal(id)
        site.name.should.equal('Gas Station')
        site.location.should.equal('Toronto, ON')
        site.tags.should.eql(['energy', 'gas'])
        done()
      }).catch(done)
    })
    it('should return a site by id', done => {
      const id = 'bbbb-bbbb'
      service.get(context, id).then(site => {
        site.id.should.equal(id)
        site.name.should.equal('Power Plant')
        site.location.should.equal('Montreal, QC')
        site.tags.should.eql(['energy', 'power'])
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should update a site by id', done => {
      const id = 'aaaa-aaaa'
      const input = { name: 'Oil Platform', location: 'Vancouver, BC' }
      service.update(context, id, input).then(site => {
        site.id.should.equal(id)
        site.name.should.equal(input.name)
        site.location.should.equal(input.location)
        done()
      }).catch(done)
    })
    it('should update a site by id', done => {
      const id = 'bbbb-bbbb'
      const input = { name: 'Hydropower Plant', location: 'Ottawa, ON' }
      service.update(context, id, input).then(site => {
        site.id.should.equal(id)
        site.name.should.equal(input.name)
        site.location.should.equal(input.location)
        done()
      }).catch(done)
    })
  })

  describe('modify', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should modify a site by id', done => {
      const id = 'aaaa-aaaa'
      const input = { priority: 2 }
      service.modify(context, id, input).then(site => {
        should.exist(site.id)
        should.exist(site.name)
        should.exist(site.location)
        should.exist(site.priority)
        should.exist(site.tags)
        done()
      }).catch(done)
    })
    it('should modify a site by id', done => {
      const id = 'bbbb-bbbb'
      const input = { priority: 4 }
      service.modify(context, id, input).then(site => {
        should.exist(site.id)
        should.exist(site.name)
        should.exist(site.location)
        should.exist(site.priority)
        should.exist(site.tags)
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    let service, context

    beforeEach(() => {
      service = new SiteService(config, { logger })
      context = {}
    })

    it('should delete a site by id', done => {
      const id = 'aaaa-aaaa'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
    it('should delete a site by id', done => {
      const id = 'bbbb-bbbb'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
  })
})
