/* eslint-env mocha */
const should = require('should')

const SwitchService = require('./switch')

describe('SwitchService', () => {
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
      const service = new SwitchService(config)
      should.exist(service.logger)
    })
    it('should create a new service with provided options', () => {
      const options = { logger }
      const service = new SwitchService(config, options)
      service.logger.should.equal(options.logger)
    })
  })

  describe('create', () => {
    let service, context

    beforeEach(() => {
      service = new SwitchService(config, { logger })
      context = {}
    })

    it('should create and persist a new switch', done => {
      const input = { siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: [ 'OFF', 'ON' ] }
      service.create(context, input).then(swtch => {
        should.exist(swtch.id)
        swtch.siteId.should.equal(input.siteId)
        swtch.name.should.equal(input.name)
        swtch.state.should.equal(input.state)
        swtch.states.should.equal(input.states)
        done()
      }).catch(done)
    })
    it('should create and persist a new switch', done => {
      const input = { siteId: 'bbbb-bbbb', name: 'Light', state: 'OFF', states: [ 'OFF', 'ON' ] }
      service.create(context, input).then(swtch => {
        should.exist(swtch.id)
        swtch.siteId.should.equal(input.siteId)
        swtch.name.should.equal(input.name)
        swtch.state.should.equal(input.state)
        swtch.states.should.equal(input.states)
        done()
      }).catch(done)
    })
  })

  describe('all', () => {
    let service, context

    beforeEach(() => {
      service = new SwitchService(config, { logger })
      context = {}
    })

    it('should return all switches of a site', done => {
      const siteId = 'aaaa-aaaa'
      service.all(context, siteId).then(switches => {
        switches.should.have.length(1)
        done()
      }).catch(done)
    })
    it('should return all switches of a site', done => {
      const siteId = 'bbbb-bbbb'
      service.all(context, siteId).then(switches => {
        switches.should.have.length(1)
        done()
      }).catch(done)
    })
  })

  describe('get', () => {
    let service, context

    beforeEach(() => {
      service = new SwitchService(config, { logger })
      context = {}
    })

    it('should return a switch by id', done => {
      const id = '3333-3333'
      service.get(context, id).then(swtch => {
        swtch.id.should.equal(id)
        swtch.siteId.should.equal('aaaa-aaaa')
        swtch.name.should.equal('Light')
        swtch.state.should.equal('OFF')
        swtch.states.should.eql(['OFF', 'ON'])
        done()
      }).catch(done)
    })
    it('should return a switch by id', done => {
      const id = '4444-4444'
      service.get(context, id).then(swtch => {
        swtch.id.should.equal(id)
        swtch.siteId.should.equal('bbbb-bbbb')
        swtch.name.should.equal('Light')
        swtch.state.should.equal('OFF')
        swtch.states.should.eql(['OFF', 'ON'])
        done()
      }).catch(done)
    })
  })

  describe('update', () => {
    let service, context

    beforeEach(() => {
      service = new SwitchService(config, { logger })
      context = {}
    })

    it('should update a switch by id', done => {
      const id = '3333-3333'
      const state = 'ON'
      service.update(context, id, { state }).then(swtch => {
        swtch.id.should.equal(id)
        swtch.siteId.should.equal('aaaa-aaaa')
        swtch.name.should.equal('Light')
        swtch.state.should.equal('ON')
        swtch.states.should.eql(['OFF', 'ON'])
        done()
      }).catch(done)
    })
    it('should update a switch by id', done => {
      const id = '4444-4444'
      const state = 'ON'
      service.update(context, id, { state }).then(swtch => {
        swtch.id.should.equal(id)
        swtch.siteId.should.equal('bbbb-bbbb')
        swtch.name.should.equal('Light')
        swtch.state.should.equal('ON')
        swtch.states.should.eql(['OFF', 'ON'])
        done()
      }).catch(done)
    })
  })

  describe('delete', () => {
    let service, context

    beforeEach(() => {
      service = new SwitchService(config, { logger })
      context = {}
    })

    it('should delete a switch by id', done => {
      const id = '3333-3333'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
    it('should delete a switch by id', done => {
      const id = '4444-4444'
      service.delete(context, id).then(() => {
        done()
      }).catch(done)
    })
  })
})
