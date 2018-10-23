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
      mongoUri: 'mongodb://mongo',
      mongoUsername: 'service',
      mongoPassword: 'password'
    }
    logger = {
      trace () {},
      debug () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    mongoose = {
      connect () {},
      disconnect () {},
      connection: new EventEmitter()
    }
    _mongoose = sinon.mock(mongoose)

    mongo = new Mongo(config, { logger, mongoose })
  })

  afterEach(() => {
    _mongoose.restore()
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
        useNewUrlParser: true,
        autoReconnect: true,
        auth: {
          user: config.mongoUsername,
          password: config.mongoPassword
        }
      }
    })

    it('should error when mongoose connect fails', done => {
      _mongoose.expects('connect').withArgs(config.mongoUri, opts).rejects(new Error('error'))
      mongo.connect().catch(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })

    it('should error when mongoose conenction fails', done => {
      _mongoose.expects('connect').withArgs(config.mongoUri, opts).resolves()
      process.nextTick(() => mongoose.connection.emit('error', new Error('error')))
      mongo.connect().catch(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })

    it('should catch close event after connection open event', done => {
      _mongoose.expects('connect').withArgs(config.mongoUri, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'open', 'close' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })

    it('should succeed when mongoose connects to Mongo', done => {
      _mongoose.expects('connect').withArgs(config.mongoUri, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'connected', 'open' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })

    it('should succeed when mongoose reconnects to Mongo after a disconnection', done => {
      _mongoose.expects('connect').withArgs(config.mongoUri, opts).resolves()
      process.nextTick(() => emitEvents(mongoose.connection, [ 'disconnected', 'connected', 'reconnected', 'open' ]))
      mongo.connect().then(conn => {
        _mongoose.verify()
        conn.should.eql(mongoose.connection)
        done()
      }).catch(done)
    })
  })

  describe('disconnect', () => {
    it('should error when mongoose disconnect fails', done => {
      _mongoose.expects('disconnect').yields(new Error('error'))
      mongo.disconnect(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should reject when mongoose disconnect fails', done => {
      _mongoose.expects('disconnect').rejects(new Error('error'))
      mongo.disconnect().catch(err => {
        _mongoose.verify()
        should.exist(err)
        err.message.should.equal('error')
        done()
      })
    })
    it('should succeed when mongoose disconnects from Mongo', done => {
      _mongoose.expects('disconnect').yields()
      mongo.disconnect(err => {
        _mongoose.verify()
        should.not.exist(err)
        done()
      })
    })
    it('should resolve when mongoose disconnects from Mongo', done => {
      _mongoose.expects('disconnect').resolves()
      mongo.disconnect().then(() => {
        _mongoose.verify()
        done()
      }).catch(done)
    })
  })
})
