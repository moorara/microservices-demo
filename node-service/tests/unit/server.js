/* eslint-env mocha */
const http = require('http')
const sinon = require('sinon')
const should = require('should')
const EventEmitter = require('events')

const Server = require('../../server')

describe('Server', () => {
  let app, routers, logger
  let mongo, _mongo
  let configProvider, _configProvider
  let httpServer
  let server

  beforeEach(() => {
    app = {
      use () {}
    }
    routers = {
      links: {
        router: {}
      }
    }
    logger = {
      debug () {},
      verbose () {},
      info () {},
      warn () {},
      error () {},
      fatal () {}
    }

    mongo = {
      connect () {}
    }
    _mongo = sinon.mock(mongo)

    configProvider = {
      getConfig () {}
    }
    _configProvider = sinon.mock(configProvider)

    httpServer = {
      close () {},
      listen (port, cb) { cb() }
    }

    server = new Server({ app, routers, logger, mongo, configProvider })
  })

  afterEach(() => {
    _mongo.restore()
    _configProvider.restore()
  })

  describe('constructor', () => {
    it('should create a new server with defaults', () => {
      server = new Server()
      should.exist(server.app)
      should.exist(server.routers)
      should.exist(server.logger)
      should.exist(server.configProvider)
    })
  })

  describe('start', () => {
    let config
    let conn
    let _http, _listen, _close

    beforeEach(() => {
      config = {
        servicePort: '10000',
        mongoUrl: 'mongodb://mongo:27017',
        mongoUser: 'user',
        mongoPass: 'pass'
      }

      conn = new EventEmitter()

      _http = sinon.stub(http, 'createServer')
      _http.callsFake(() => httpServer)

      _listen = sinon.spy(httpServer, 'listen')
      _close = sinon.spy(httpServer, 'close')
    })

    afterEach(() => {
      _http.restore()
      _listen.restore()
      _close.restore()
    })

    it('should error when cannot get config', done => {
      _configProvider.expects('getConfig').withArgs().rejects(new Error('error'))
      server.start().catch(err => {
        should.exist(err)
        err.message.should.equal('error')
        _configProvider.verify()
        done()
      })
    })
    it('should error when cannot connect to mongo', done => {
      _configProvider.expects('getConfig').withArgs().resolves(config)
      _mongo.expects('connect').withArgs().rejects(new Error('error'))
      server.start().catch(err => {
        should.exist(err)
        err.message.should.equal('error')
        _configProvider.verify()
        _mongo.verify()
        done()
      })
    })
    it('should close server when mongo connection fails', done => {
      _configProvider.expects('getConfig').withArgs().resolves(config)
      _mongo.expects('connect').withArgs().resolves(conn)
      server.start().then(() => {
        conn.emit('error', new Error('error'))
        _configProvider.verify()
        _mongo.verify()
        _listen.calledWith(config.servicePort).should.be.true()
        _close.called.should.be.true()
        done()
      }).catch(done)
    })
    it('should listen on given port successfully', done => {
      _configProvider.expects('getConfig').withArgs().resolves(config)
      _mongo.expects('connect').withArgs().resolves(conn)
      server.start().then(() => {
        _configProvider.verify()
        _mongo.verify()
        sinon.assert.calledWith(httpServer.listen, config.servicePort)
        done()
      }).catch(done)
    })
  })
})
