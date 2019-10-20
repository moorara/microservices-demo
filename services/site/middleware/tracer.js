const onFinished = require('on-finished')
const opentracing = require('opentracing')

const { createTracer } = require('../util/tracer')

module.exports.create = options => {
  options = options || {}
  options.tracer = options.tracer || createTracer({ serviceName: 'middleware' })

  return (req, res, next) => {
    // https://opentracing-javascript.surge.sh
    const parentSpanContext = options.tracer.extract(opentracing.FORMAT_HTTP_HEADERS, req.headers)
    const span = options.tracer.startSpan('http-request', { childOf: parentSpanContext })
    req.context = { span }

    onFinished(res, (err, res) => {
      if (err) {
        span.log({ message: err.message })
      }

      // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
      span.setTag('http.version', req.httpVersion)
      span.setTag(opentracing.Tags.HTTP_METHOD, req.method)
      span.setTag(opentracing.Tags.HTTP_URL, req.originalUrl)
      span.setTag(opentracing.Tags.HTTP_STATUS_CODE, res.statusCode)
      span.finish()
    })

    next()
  }
}
