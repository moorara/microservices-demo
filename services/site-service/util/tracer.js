const { initTracer } = require('jaeger-client')

const Logger = require('./logger')

module.exports.createTracer = (config, options) => {
  options = options || {}
  options.logger = options.logger || new Logger('Tracer')
  options.metrics = options.metrics

  // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/configuration.js
  const tracerConfig = {
    serviceName: config.serviceName,
    sampler: {
      type: 'const',
      param: 1
    },
    reporter: {
      agentHost: config.jaegerAgentHost,
      agentPort: config.jaegerAgentPort,
      logSpans: config.jaegerReporterLogSpans
    }
  }

  const tracerOptions = {
    logger: options.logger, // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/logger.js
    metrics: options.metrics // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/metrics.js
  }

  return initTracer(tracerConfig, tracerOptions)
}
