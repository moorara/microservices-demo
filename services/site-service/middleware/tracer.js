const onFinished = require('on-finished')

const { createTracer } = require('../util/tracer')

module.exports.create = options => {
  options = options || {}
  options.tracer = options.tracer || createTracer('site-service-middleware')

  return (req, res, next) => {
    const span = options.tracer.startSpan('http-request')
    req.context = { span }

    onFinished(res, (err, res) => {
      if (err) {
        return
      }

      // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
      span.setTag('http.method', req.method)
      span.setTag('http.url', req.originalUrl)
      span.setTag('http.status_code', res.statusCode)
      span.log({
        event: 'http-request',
        method: req.method,
        url: req.originalUrl,
        statusCode: res.statusCode
      })

      span.finish()
    })

    next()
  }
}
