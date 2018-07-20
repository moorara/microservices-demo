const onFinished = require('on-finished')

const { createTracer } = require('../util/tracer')

module.exports.http = options => {
  options = options || {}
  options.tracer = options.tracer || createTracer('site-service-middleware')

  return (req, res, next) => {
    const span = options.tracer.startSpan('http-request')

    onFinished(res, (err, res) => {
      if (err) {
        return
      }

      span.setTag('http.method', req.method)
      span.setTag('http.url', req.endpoint)
      span.setTag('http.status_code', res.statusCode)
      span.log({
        event: 'http-request',
        method: req.method,
        url: req.endpoint,
        statusCode: res.statusCode
      })
      span.finish()
    })

    next()
  }
}
