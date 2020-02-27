/* eslint-env mocha */
const http = require('http')
const EventEmitter = require('events')
const promClient = require('prom-client')
const opentracing = require('opentracing')
const sinon = require('sinon')
const should = require('should')

const Server = require('../../server')

describe('Server', () => {
  describe('constructor', () => {
    afterEach(() => {
      // Clear default metrics register, so subsequent tests will not fail.
      promClient.register.clear()
    })

    it('should create a new server with defaults', () => {
      const server = new Server()
      should.exist(server.app)
      should.exist(server.routers)
      should.exist(server.logger)
      should.exist(server.metrics)
      should.exist(server.configProvider)
    })
  })

  describe('start', () => {
    let httpServer, _http, _listen, _close
    let app, routers
    let mongo, _mongo
    let configProvider, _configProvider
    let logger, metrics, tracer
    let server
    let config, conn

    beforeEach(() => {
      httpServer = {
        close () {},
        listen (port, cb) { cb() }
      }
      _http = sinon.stub(http, 'createServer')
      _http.callsFake(() => httpServer)
      _listen = sinon.spy(httpServer, 'listen')
      _close = sinon.spy(httpServer, 'close')

      app = {
        use () {}
      }
      routers = {
        sites: {
          router: {}
        }
      }

      mongo = {
        connect () {},
        disconnect () {}
      }
      _mongo = sinon.mock(mongo)

      configProvider = {
        getConfig () {}
      }
      _configProvider = sinon.mock(configProvider)

      logger = {
        trace () {},
        debug () {},
        info () {},
        warn () {},
        error () {},
        fatal () {}
      }
      metrics = {
        router: {},
        register: new promClient.Registry()
      }
      tracer = new opentracing.MockTracer()

      server = new Server({ app, routers, mongo, configProvider, logger, metrics, tracer })

      config = {
        servicePort: '10000',
        mongoUri: 'mongodb://mongo:27017',
        mongoUsername: 'user',
        mongoPassword: 'pass',
        jaegerAgentHost: 'jaeger-agent',
        jaegerAgentPort: 6832,
        jaegerLogSpans: true
      }
      conn = new EventEmitter()
    })

    afterEach(() => {
      _http.restore()
      _listen.restore()
      _close.restore()
      _mongo.restore()
      _configProvider.restore()
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
