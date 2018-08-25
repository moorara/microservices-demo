const express = require('express')
const promClient = require('prom-client')
const opentracing = require('opentracing')
const expressWinston = require('express-winston')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

const histogramName = 'http_requests_duration_seconds'
const histogramHelp = 'duration histogram of http requests'
const summaryName = 'http_requests_duration_quantiles_seconds'
const summaryHelp = 'duration summary of http requests'
const labelNames = [ 'httpVersion', 'method', 'url', 'statusCode' ]
const buckets = [ 0.01, 0.1, 0.5, 1 ]
const percentiles = [ 0.1, 0.5, 0.95, 0.99 ]
const defaultSpanName = 'http-request'

class MonitorMiddleware {
  /**
   * Creates a middleware for monitoring http requests and responses
   * @param {object} options          Optional
   * @param {object} options.winston  A Winston logger instance
   * @param {object} options.metadata Metadata for logger
   * @param {object} options.register Prometheus registry for metrics
   * @param {object} options.tracer   An OpenTracing instance
   * @param {string} options.spanName The name of span
   */
  constructor (options) {
    options = options || {}
    const winston = options.winston || Logger.getWinstonLogger()
    const metadata = Object.assign({}, Logger.metadata, options.metadata, { module: 'MonitorMiddleware' })
    const register = options.register || promClient.register
    this.tracer = options.tracer || createTracer({ serviceName: 'middleware' })
    this.spanName = options.spanName || defaultSpanName

    const loggerMiddleware = expressWinston.logger({
      statusLevels: true,
      meta: true,
      baseMeta: metadata,
      winstonInstance: winston,
      expressFormat: process.env.NODE_ENV === 'development'
    })

    this.histogram = new promClient.Histogram({
      name: histogramName,
      help: histogramHelp,
      labelNames,
      buckets,
      registers: [ register ]
    })

    this.summary = new promClient.Summary({
      name: summaryName,
      help: summaryHelp,
      labelNames,
      percentiles,
      registers: [ register ]
    })

    this.router = express.Router()
    this.router.use(loggerMiddleware)
    this.router.use(this._middleware.bind(this))
  }

  _middleware (req, res, next) {
    const startTime = Date.now()
    const span = this.tracer.startSpan(this.spanName)
    req.context = { span }

    const end = res.end
    res.end = (data, encoding) => {
      res.end = end
      res.end(data, encoding)

      const duration = (Date.now() - startTime) / 1000
      const url = req.path

      // Metrics
      const labelValues = {
        httpVersion: req.httpVersion,
        method: req.method,
        url: url,
        statusCode: res.statusCode
      }
      this.histogram.observe(labelValues, duration)
      this.summary.observe(labelValues, duration)

      // Traces
      // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
      span.setTag('http.version', req.httpVersion)
      span.setTag(opentracing.Tags.HTTP_METHOD, req.method)
      span.setTag(opentracing.Tags.HTTP_URL, url)
      span.setTag(opentracing.Tags.HTTP_STATUS_CODE, res.statusCode)
      span.finish()
    }

    next()
  }
}

module.exports = MonitorMiddleware
