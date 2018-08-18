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
      should.exist(server.configProvider)
      should.exist(server.logger)
      should.exist(server.metrics)
      should.exist(server.app)
      should.exist(server.routers)
      should.exist(server.middleware)
    })
  })

  describe('start', () => {
    let httpServer, _createServer, _listen, _close
    let config, configProvider, _configProvider
    let logger, metrics, tracer
    let app, routers, middleware
    let server

    beforeEach(() => {
      httpServer = {
        close () {},
        listen (port, callback) { callback() }
      }
      _createServer = sinon.stub(http, 'createServer')
      _createServer.callsFake(() => httpServer)
      _listen = sinon.spy(httpServer, 'listen')
      _close = sinon.spy(httpServer, 'close')

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

      app = {
        use () {}
      }
      routers = {
        liveness: (req, res) => res.sendStatus(200),
        readiness: (req, res) => res.sendStatus(200),
        graphql: {
          router: (req, res) => res.sendStatus(200)
        }
      }
      middleware = {
        monitor: {
          router: (req, res, next) => next()
        }
      }

      server = new Server({
        configProvider,
        logger,
        metrics,
        tracer,
        app,
        routers,
        middleware
      })
    })

    afterEach(() => {
      _configProvider.restore()
      _createServer.restore()
      _listen.restore()
      _close.restore()
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
