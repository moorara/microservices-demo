/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')

const resolvers = require('./switch')

describe('switchResolvers', () => {
  let span, logger
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

    switchService = {
      get () {},
      all () {},
      create () {},
      update () {},
      delete () {}
    }
    _switchService = sinon.mock(switchService)

    context = { span, logger, switchService }
    info = {}
  })

  afterEach(() => {
    _switchService.restore()
  })

  describe('Query', () => {
    describe('switch', () => {
      let id

      it('should throw an error when service.get fails', done => {
        id = '3333-3333'
        const err = new Error('get error')
        _switchService.expects('get').withArgs({ span }, id).rejects(err)
        resolvers.Query.switch(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves to a switch', done => {
        id = '4444-4444'
        const swtch = { id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }
        _switchService.expects('get').withArgs({ span }, id).resolves(swtch)
        resolvers.Query.switch(null, { id }, context, info).then(s => {
          s.should.eql(swtch)
          _switchService.verify()
          done()
        }).catch(done)
      })
    })

    describe('switches', () => {
      let siteId

      it('should throw an error when service.all fails', done => {
        siteId = 'aaaa-aaaa'
        const err = new Error('all error')
        _switchService.expects('all').withArgs({ span }, siteId).rejects(err)
        resolvers.Query.switches(null, { siteId }, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves to an array switches', done => {
        siteId = 'bbbb-bbbb'
        const switches = [{ id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }]
        _switchService.expects('all').withArgs({ span }, siteId).resolves(switches)
        resolvers.Query.switches(null, { siteId }, context, info).then(s => {
          s.should.eql(switches)
          _switchService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Mutation', () => {
    describe('installSwitch', () => {
      let input

      it('should throw an error when service.create fails', done => {
        input = {}
        const err = new Error('create error')
        _switchService.expects('create').withArgs({ span }, input).rejects(err)
        resolvers.Mutation.installSwitch(null, { input }, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the new switch', done => {
        input = { siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }
        const swtch = Object.assign({}, { id: '4444-4444' }, input)
        _switchService.expects('create').withArgs({ span }, input).resolves(swtch)
        resolvers.Mutation.installSwitch(null, { input }, context, info).then(s => {
          s.should.eql(swtch)
          _switchService.verify()
          done()
        }).catch(done)
      })
    })

    describe('setSwitch', () => {
      let id, state

      it('should throw an error when service.update fails', done => {
        id = '3333-3333'
        state = ''
        const err = new Error('update error')
        _switchService.expects('update').withArgs({ span }, id, { state }).rejects(err)
        resolvers.Mutation.setSwitch(null, { id, state }, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves to the updated switch', done => {
        id = '4444-4444'
        state = 'ON'
        const swtch = { siteId: 'aaaa-aaaa', name: 'Light', state: 'ON', states: ['OFF', 'ON'] }
        _switchService.expects('update').withArgs({ span }, id, { state }).resolves(swtch)
        resolvers.Mutation.setSwitch(null, { id, state }, context, info).then(s => {
          s.should.eql(swtch)
          _switchService.verify()
          done()
        }).catch(done)
      })
    })

    describe('removeSwitch', () => {
      let id

      it('should throw an error when service.delete fails', done => {
        id = '3333-3333'
        const err = new Error('delete error')
        _switchService.expects('delete').withArgs({ span }, id).rejects(err)
        resolvers.Mutation.removeSwitch(null, { id }, context, info).catch(e => {
          e.should.eql(err)
          _switchService.verify()
          done()
        })
      })
      it('should return a promise that resolves true', done => {
        id = '4444-4444'
        _switchService.expects('delete').withArgs({ span }, id).resolves(true)
        resolvers.Mutation.removeSwitch(null, { id }, context, info).then(result => {
          result.should.be.true()
          _switchService.verify()
          done()
        }).catch(done)
      })
    })
  })

  describe('Switch', () => {
    let swtch

    beforeEach(() => {
      swtch = { id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }
    })

    describe('id', () => {
      it('should return switch id', () => {
        const id = resolvers.Switch.id(swtch)
        should.equal(id, swtch.id)
      })
    })
    describe('siteId', () => {
      it('should return switch siteId', () => {
        const siteId = resolvers.Switch.siteId(swtch)
        should.equal(siteId, swtch.siteId)
      })
    })
    describe('name', () => {
      it('should return switch name', () => {
        const name = resolvers.Switch.name(swtch)
        should.equal(name, swtch.name)
      })
    })
    describe('state', () => {
      it('should return switch state', () => {
        const state = resolvers.Switch.state(swtch)
        should.equal(state, swtch.state)
      })
    })
    describe('states', () => {
      it('should return switch states', () => {
        const states = resolvers.Switch.states(swtch)
        should.equal(states, swtch.states)
      })
    })
  })
})
