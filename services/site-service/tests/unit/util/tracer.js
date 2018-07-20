/* eslint-env mocha */
const should = require('should')

const { createTracer } = require('../../../util/tracer')

describe('Tracer', () => {
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
      tracer = createTracer(config, options)

      should.exist(tracer)
      tracer.close(done)
    })
    it('should create a new tracer with non-default agent', done => {
      const config = {
        serviceName: 'node-service',
        jaegerAgentHost: 'jaeger-agent',
        jaegerAgentPort: 6831
      }
      tracer = createTracer(config, options)

      should.exist(tracer)
      tracer.close(done)
    })
    it('should create a new tracer with logging spans enabled', done => {
      const config = {
        serviceName: 'node-service',
        jaegerReporterLogSpans: true
      }
      tracer = createTracer(config, options)

      should.exist(tracer)
      tracer.close(done)
    })
  })
})
