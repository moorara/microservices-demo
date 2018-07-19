const promClient = require('prom-client')
const onFinished = require('on-finished')

const histogramName = 'http_requests_duration_seconds'
const summaryName = 'http_requests_duration_quantiles_seconds'
const defaultLabels = [ 'method', 'endpoint', 'statusCode', 'statusClass' ]
const defaultBuckets = [ 0.01, 0.1, 0.5, 1 ]
const defaultPercentiles = [ 0.1, 0.5, 0.95, 0.99 ]

module.exports.create = options => {
  options = options || {}
  options.register = options.register || promClient.register

  const httpHistogram = new promClient.Histogram({
    name: histogramName,
    help: 'duration histogram of http requests',
    labelNames: defaultLabels,
    buckets: defaultBuckets,
    registers: [ options.register ]
  })

  const httpSummary = new promClient.Summary({
    name: summaryName,
    help: 'duration summary of http requests',
    labelNames: defaultLabels,
    percentiles: defaultPercentiles,
    registers: [ options.register ]
  })

  return (req, res, next) => {
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

      httpHistogram.observe(labelValues, duration)
      httpSummary.observe(labelValues, duration)
    })

    next()
  }
}
