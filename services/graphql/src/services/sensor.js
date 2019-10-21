const http = require('http')
const axios = require('axios')
const opentracing = require('opentracing')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

const timeout = 1000

class SensorService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SensorService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'SensorService' })
    this.axios = options.axios || axios.create({
      timeout,
      httpAgent: new http.Agent({ keepAlive: true }),
      baseURL: `http://${config.sensorServiceAddr}/v1/`
    })
  }

  async exec (context, name, func) {
    let err, result
    let startTime, latency

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.SPAN_KIND, 'client')
    span.setTag(opentracing.Tags.PEER_SERVICE, 'sensor-service')
    span.setTag(opentracing.Tags.PEER_ADDRESS, this.axios.defaults.baseUrl)

    const headers = {}
    this.tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers)

    // Core functionality
    try {
      startTime = Date.now()
      result = await func(headers)
    } catch (e) {
      err = e
      this.logger.error(err)
    } finally {
      latency = (Date.now() - startTime) / 1000
    }

    // Metrics
    const labelValues = { op: name, success: err ? 'false' : 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Tracing
    span.log({
      event: name,
      message: err ? err.message : 'successful!'
    })
    span.finish()

    if (err) {
      throw err
    }
    return result
  }

  create (context, input) {
    return this.exec(context, 'create-sensor', headers => {
      return this.axios.request({
        headers,
        method: 'post',
        url: '/sensors',
        data: input
      }).then(res => res.data)
    })
  }

  all (context, siteId) {
    return this.exec(context, 'get-sensors', headers => {
      return this.axios.request({
        headers,
        method: 'get',
        url: '/sensors',
        params: { siteId }
      }).then(res => res.data)
    })
  }

  get (context, id) {
    return this.exec(context, 'get-sensor', headers => {
      return this.axios.request({
        headers,
        method: 'get',
        url: `/sensors/${id}`
      }).then(res => res.data)
    })
  }

  update (context, id, input) {
    return this.exec(context, 'update-sensor', async headers => {
      await this.axios.request({
        headers,
        method: 'put',
        url: `/sensors/${id}`,
        data: input
      })

      // sensor-service does not respond with updated sensor
      const updated = await this.axios.request({
        headers,
        method: 'get',
        url: `/sensors/${id}`
      })

      return updated.data
    })
  }

  delete (context, id) {
    return this.exec(context, 'delete-sensor', headers => {
      return this.axios.request({
        headers,
        method: 'delete',
        url: `/sensors/${id}`
      }).then(res => res.data)
    })
  }
}

module.exports = SensorService
