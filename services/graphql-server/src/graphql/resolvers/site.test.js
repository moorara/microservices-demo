/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')

const resolvers = require('./site')

describe('siteResolvers', () => {
  let span, logger
  let siteService, _siteService
  let sensorService, _sensorService
  let switchService, _switchService
  let context, info

  beforeEach(() => {
    span = {}
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    siteService = {
      get () {},
      all () {},
      create () {},
      update () {},
      delete () {}
    }
    _siteService = sinon.mock(siteService)

    sensorService = {
      all () {}
    }
    _sensorService = sinon.mock(sensorService)

    switchService = {
      getSwitches () {}
    }
    _switchService = sinon.mock(switchService)

    context = { span, logger, siteService, sensorService, switchService }
    info = {}
  })

  afterEach(() => {
    _siteService.restore()
    _sensorService.restore()
    _switchService.restore()
  })

  describe('Query', () => {
    describe('site', () => {
      let id

      it('should throw an error when service.get fails', done => {
        id = 'aaaa-aaaa'
        const err = new Error('get error')
        _siteService.expects('get').withArgs({ span }, id).rejects(err)
        resolvers.Query.site(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _siteService.verify()
          done()
        })
      })
      it('should return a promise that resolves to a site', done => {
        id = 'bbbb-bbbb'
        const site = { id: 'aaaa-aaaa', name: 'Power Plant', location: 'Ottawa, ON', tags: ['energy', 'power'] }
        _siteService.expects('get').withArgs({ span }, id).resolves(site)
        resolvers.Query.site(null, { id }, context, info).then(s => {
          s.should.eql(site)
          _siteService.verify()
          done()
        }).catch(done)
      })
    })

    describe('sites', () => {
      let args // name, location, tags, minPriority, maxPriority, limit, skip

      it('should throw an error when service.all fails', done => {
        const err = new Error('all error')
        _siteService.expects('all').withArgs({ span }, args).rejects(err)
        resolvers.Query.sites(null, args, context, info).catch(e => {
          e.should.eql(err)
          _siteService.verify()
          done()
        })
      })
      it('should return a promise that resolves to an array sites', done => {
        const sites = [{ id: 'aaaa-aaaa', name: 'Power Plant', location: 'Ottawa, ON', tags: ['energy', 'power'] }]
        _siteService.expects('all').withArgs({ span }, args).resolves(sites)
        resolvers.Query.sites(null, args, context, info).then(s => {
          s.should.eql(sites)
          _siteService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Mutation', () => {
    describe('createSite', () => {
      let input

      it('should throw an error when service.create fails', done => {
        input = {}
        const err = new Error('create error')
        _siteService.expects('create').withArgs({ span }, input).rejects(err)
        resolvers.Mutation.createSite(null, { input }, context, info).catch(e => {
          e.should.eql(err)
          _siteService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the new site', done => {
        input = { name: 'Power Plant', location: 'Ottawa, ON', tags: ['energy', 'power'] }
        const site = Object.assign({}, { id: 'bbbb-bbbb' }, input)
        _siteService.expects('create').withArgs({ span }, input).resolves(site)
        resolvers.Mutation.createSite(null, { input }, context, info).then(s => {
          s.should.eql(site)
          _siteService.verify()
          done()
        }).catch(done)
      })
    })

    describe('updateSite', () => {
      let id, input

      it('should throw an error when service.update fails', done => {
        id = 'aaaa-aaaa'
        input = {}
        const err = new Error('update error')
        _siteService.expects('update').withArgs({ span }, id, input).rejects(err)
        resolvers.Mutation.updateSite(null, { id, input }, context, info).catch(e => {
          e.should.eql(err)
          _siteService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the updated site', done => {
        id = 'bbbb-bbbb'
        input = { name: 'Power Plant', location: 'Ottawa, ON', tags: ['energy', 'power'] }
        const site = Object.assign({}, { id }, input)
        _siteService.expects('update').withArgs({ span }, id, input).resolves(site)
        resolvers.Mutation.updateSite(null, { id, input }, context, info).then(s => {
          s.should.eql(site)
          _siteService.verify()
          done()
        }).catch(done)
      })
    })

    describe('deleteSite', () => {
      let id

      it('should throw an error when service.delete fails', done => {
        id = 'aaaa-aaaa'
        const err = new Error('delete error')
        _siteService.expects('delete').withArgs({ span }, id).rejects(err)
        resolvers.Mutation.deleteSite(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _siteService.verify()
          done()
        })
      })
      it('should return a promise that resolves true', done => {
        id = 'bbbb-bbbb'
        _siteService.expects('delete').withArgs({ span }, id).resolves(true)
        resolvers.Mutation.deleteSite(null, { id }, context, info).then(result => {
          result.should.be.true()
          _siteService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Site', () => {
    let site

    beforeEach(() => {
      site = { id: 'aaaa-aaaa', name: 'Power Plant', location: 'Ottawa, ON', tags: ['energy', 'power'] }
    })

    describe('id', () => {
      it('should return site id', () => {
        const id = resolvers.Site.id(site)
        should.equal(id, site.id)
      })
    })
    describe('name', () => {
      it('should return site name', () => {
        const name = resolvers.Site.name(site)
        should.equal(name, site.name)
      })
    })
    describe('location', () => {
      it('should return site location', () => {
        const location = resolvers.Site.location(site)
        should.equal(location, site.location)
      })
    })
    describe('priority', () => {
      it('should return site priority', () => {
        const priority = resolvers.Site.priority(site)
        should.equal(priority, site.priority)
      })
    })
    describe('tags', () => {
      it('should return site tags', () => {
        const tags = resolvers.Site.tags(site)
        should.equal(tags, site.tags)
      })
    })

    describe('sensors', () => {
      it('should return a promise that rejects with an error when service.all fails', done => {
        const err = new Error('all error')
        _sensorService.expects('all').withArgs({ span }, site.id).rejects(err)
        resolvers.Site.sensors(site, {}, context, info).catch(e => {
          e.should.eql(err)
          _sensorService.verify()
          done()
        })
      })
      it('should return a promise that resolves to site sensors', done => {
        const sensors = [{ id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 }]
        _sensorService.expects('all').withArgs({ span }, site.id).resolves(sensors)
        resolvers.Site.sensors(site, {}, context, info).then(s => {
          s.should.eql(sensors)
          _sensorService.verify()
          done()
        }).catch(done)
      })
    })

    describe('switches', () => {
      it('should return a promise that rejects with an error when service.all fails', done => {
        const err = new Error('all error')
        _switchService.expects('getSwitches').withArgs({ span }, site.id).rejects(err)
        resolvers.Site.switches(site, {}, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves to site switches', done => {
        const switches = [{ id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }]
        _switchService.expects('getSwitches').withArgs({ span }, site.id).resolves(switches)
        resolvers.Site.switches(site, {}, context, info).then(s => {
          s.should.eql(switches)
          _switchService.verify()
          done()
        }).catch(done)
      })
    })
  })
})
