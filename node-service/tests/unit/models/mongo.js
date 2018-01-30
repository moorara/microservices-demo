/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')
const EventEmitter = require('events')

const Mongo = require('../../../models/mongo')

const emitEvents = (emitter, events) => {
  if (events.length > 0) {
    process.nextTick(() => {
      let event = events.shift()
      emitter.emit(event)
      emitEvents(emitter, events)
    })
  }
}

describe('Mongo', () => {
  let config, logger
  let mongoose, _mongoose
  let mongo

  beforeEach(() => {
    config = {
      mongoUrl: 'mongodb://mongo',
      mongoUser: 'user',
      mongoPass: 'pass'
    }
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    mongoose = {
      connect () {},
      connection: new EventEmitter()
    }
    _mongoose = sinon.mock(mongoose)

    mongo = new Mongo(config, {
      logger,
      mongoose
    })
  })

  describe('constructor', () => {
    it('should create a new instance with defaults', () => {
      mongo = new Mongo({})
      should.exist(mongo.logger)
      should.exist(mongo.mongoose)
    })
  })

  describe('connect', () => {
    let opts

    beforeEach(() => {
      opts = {
        autoReconnect: true,
        auth: {
          user: config.mongoUser,
          pass: config.mongoPass
        }
      }
    })

    afterEach(() => {
      _mongoose.restore()
    })

    it('should error when mongoose connect fails', done => {
      _mongoose.expects('connect').withArgs(config.mongoUrl, opts).rejects(new Error('error'))
      mongo.connect().catch(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })

    it('should error when mongoose conenction fails', done => {
      _mongoose.expects('connect').withArgs(config.mongoUrl, opts).resolves()
      process.nextTick(() => mongoose.connection.emit('error', new Error('error')))
      mongo.connect().catch(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })

    it('should catch close event after connection open event', done => {
      _mongoose.expects('connect').withArgs(config.mongoUrl, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'open', 'close' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })

    it('should succeed when mongoose connects to Mongo', done => {
      _mongoose.expects('connect').withArgs(config.mongoUrl, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'connected', 'open' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })

    it('should succeed when mongoose reconnects to Mongo after a disconnection', done => {
      _mongoose.expects('connect').withArgs(config.mongoUrl, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'disconnected', 'connected', 'reconnected', 'open' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })
  })
})
