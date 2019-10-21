/* eslint-env mocha */
const sinon = require('sinon')
const should = require('should')

const { MetricsFactory, createTracer } = require('./tracer')

describe('Tracer', () => {
  describe('MetricsFactory', () => {
    describe('constructor', () => {
      let metricsFactory

      it('should create a new metrics factory', () => {
        metricsFactory = new MetricsFactory('service-name')
        should.exist(metricsFactory.factory)
        metricsFactory.namespace.should.equal('service_name')
      })
    })

    describe('sanitizeName', () => {
      let factory, metricsFactory

      beforeEach(() => {
        factory = {}
        metricsFactory = new MetricsFactory('service', { factory })
      })

      it('should contain no invalid character', () => {
        const name = metricsFactory.sanitizeName('traces')
        name.should.equal('traces')
      })
      it('should contain no invalid character', () => {
        const name = metricsFactory.sanitizeName('jaeger-traces')
        name.should.equal('jaeger_traces')
      })
      it('should contain no invalid character', () => {
        const name = metricsFactory.sanitizeName('jaeger:traces')
        name.should.equal('jaeger_traces')
      })
      it('should contain no invalid character', () => {
        const name = metricsFactory.sanitizeName('jaeger:traces-started')
        name.should.equal('jaeger_traces_started')
      })
    })

    describe('createCounter', () => {
      let factory, _factory
      let metricsFactory

      beforeEach(() => {
        factory = { createCounter () {} }
        _factory = sinon.mock(factory)
        metricsFactory = new MetricsFactory('service', { factory })
      })

      afterEach(() => {
        _factory.verify()
        _factory.restore()
      })

      it('should sanitize name before passing it to factory', () => {
        _factory.expects('createCounter').withArgs('service_metric_name').returns()
        metricsFactory.createCounter('service:metric-name')
      })
    })

    describe('createGauge', () => {
      let factory, _factory
      let metricsFactory

      beforeEach(() => {
        factory = { createGauge () {} }
        _factory = sinon.mock(factory)
        metricsFactory = new MetricsFactory('service', { factory })
      })

      afterEach(() => {
        _factory.verify()
        _factory.restore()
      })

      it('should sanitize name before passing it to factory', () => {
        _factory.expects('createGauge').withArgs('service_metric_name').returns()
        metricsFactory.createGauge('service:metric-name')
      })
    })

    describe('createTimer', () => {
      let factory, _factory
      let metricsFactory

      beforeEach(() => {
        factory = { createTimer () {} }
        _factory = sinon.mock(factory)
        metricsFactory = new MetricsFactory('service', { factory })
      })

      afterEach(() => {
        _factory.verify()
        _factory.restore()
      })

      it('should sanitize name before passing it to factory', () => {
        _factory.expects('createTimer').withArgs('service_metric_name').returns()
        metricsFactory.createTimer('service:metric-name')
      })
    })
  })

  describe('createTracer', () => {
    let tracer
    let options

    beforeEach(() => {
      options = {
        // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/logger.js
        logger: {
          info () {},
          error () {}
        },
        // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/metrics.js
        metrics: {
          createCounter: () => ({ increment () {} }),
          createTimer: () => ({ record () {} }),
          createGauge: () => ({ update () {} })
        }
      }
    })

    it('should create a new tracer with defaults', done => {
      const config = {
        serviceName: 'node-service'
      }
      tracer = createTracer(config)

      should.exist(tracer)
      tracer.close(done)
    })
    it('should create a new tracer with provided options', done => {
      const config = {
        serviceName: 'node-service',
        jaegerAgentHost: 'jaeger-agent',
        jaegerAgentPort: 6832,
        jaegerLogSpans: true
      }
      tracer = createTracer(config, options)

      should.exist(tracer)
      tracer.close(done)
    })
  })
})
