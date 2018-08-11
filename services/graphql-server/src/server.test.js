/* eslint-env mocha */
const http = require('http')
const sinon = require('sinon')
const should = require('should')
const promClient = require('prom-client')
const opentracing = require('opentracing')

const Server = require('./server')

describe('Server', () => {
  describe('constructor', () => {
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
    let httpServer, _createServer, _listen, _close
    let app, routers, logger, metrics, tracer
    let config, configProvider, _configProvider
    let server

    beforeEach(() => {
      httpServer = {
        close () {},
        listen (port, cb) { cb() }
      }
      _createServer = sinon.stub(http, 'createServer')
      _createServer.callsFake(() => httpServer)
      _listen = sinon.spy(httpServer, 'listen')
      _close = sinon.spy(httpServer, 'close')

      app = {
        use () {}
      }
      routers = {}

      logger = {
        debug () {},
        verbose () {},
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

      config = {
        servicePort: '10000',
        jaegerAgentHost: 'jaeger-agent',
        jaegerAgentPort: 6832,
        jaegerLogSpans: true
      }
      configProvider = {
        getConfig () {}
      }
      _configProvider = sinon.mock(configProvider)

      server = new Server({
        app,
        routers,
        logger,
        metrics,
        tracer,
        configProvider
      })
    })

    afterEach(() => {
      _createServer.restore()
      _listen.restore()
      _close.restore()
      _configProvider.restore()
    })

    it('should error when cannot get config', done => {
      _configProvider.expects('getConfig').rejects(new Error('config error'))
      server.start().catch(err => {
        should.exist(err)
        err.message.should.equal('config error')
        _configProvider.verify()
        done()
      })
    })
    it('should listen on service port successfully', done => {
      _configProvider.expects('getConfig').resolves(config)
      server.start().then(() => {
        _configProvider.verify()
        sinon.assert.calledWith(httpServer.listen, config.servicePort)
        done()
      }).catch(done)
    })
  })
})
