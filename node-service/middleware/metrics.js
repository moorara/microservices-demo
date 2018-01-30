const express = require('express')
const promClient = require('prom-client')
const onFinished = require('on-finished')

const histogramName = 'http_requests_duration_seconds'
const summaryName = 'http_requests_duration_quantiles_seconds'
const defaultLabels = [ 'method', 'endpoint', 'statusCode', 'statusClass' ]
const defaultBuckets = [ 0.01, 0.1, 0.5, 1 ]
const defaultPercentiles = [ 0.1, 0.5, 0.95, 0.99 ]

class MetricsMiddleware {
  constructor (options) {
    options = options || {}

    this.registery = new promClient.Registry()
    this.httpHistogram = new promClient.Histogram({
      name: histogramName,
      help: 'duration histogram of http requests',
      labelNames: defaultLabels,
      buckets: defaultBuckets,
      registers: [ this.registery ]
    })
    this.httpSummary = new promClient.Summary({
      name: summaryName,
      help: 'duration summary of http requests',
      labelNames: defaultLabels,
      percentiles: defaultPercentiles,
      registers: [ this.registery ]
    })

    this.interval = promClient.collectDefaultMetrics({
      register: this.registery
    })

    this.router = express.Router()
    this.router.get('/metrics', this.getMetrics.bind(this))
    this.router.use(this.observeDuration.bind(this))
  }

  getMetrics (req, res) {
    res.type('text')
    res.send(this.registery.metrics())
  }

  observeDuration (req, res, next) {
    let startTime = +new Date()

    onFinished(res, (err, res) => {
      if (err) {
        return
      }

      let endTime = +new Date()
      const duration = (endTime - startTime) / 1000

      const labelValues = {
        method: req.method,
        endpoint: req.endpoint,
        statusCode: res.statusCode,
        statusClass: res.statusClass
      }

      this.httpHistogram.observe(labelValues, duration)
      this.httpSummary.observe(labelValues, duration)
    })

    next()
  }
}

module.exports = MetricsMiddleware
