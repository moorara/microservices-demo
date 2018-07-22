const promClient = require('prom-client')
const { initTracer, PrometheusMetricsFactory } = require('jaeger-client')

const Logger = require('./logger')

class MetricsFactory {
  constructor (serviceName, options) {
    options = options || {}
    this.namespace = this.sanitizeName(serviceName)
    this.factory = options.factory || new PrometheusMetricsFactory(promClient, this.namespace)
  }

  sanitizeName (name) {
    return name.replace(/[-:]/g, '_')
  }

  createCounter (name, tags) {
    name = this.sanitizeName(name)
    return this.factory.createCounter(name, tags)
  }

  createGauge (name, tags) {
    name = this.sanitizeName(name)
    return this.factory.createGauge(name, tags)
  }

  createTimer (name, tags) {
    name = this.sanitizeName(name)
    return this.factory.createTimer(name, tags)
  }
}

function createTracer (config, options) {
  options = options || {}
  options.logger = options.logger || new Logger('Tracer')
  options.metrics = options.metrics || new MetricsFactory(config.serviceName)

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
      logSpans: config.jaegerLogSpans
    }
  }

  const tracerOptions = {
    logger: options.logger, // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/logger.js
    metrics: options.metrics // https://github.com/jaegertracing/jaeger-client-node/blob/master/src/_flow/metrics.js
  }

  return initTracer(tracerConfig, tracerOptions)
}

module.exports = { MetricsFactory, createTracer }
