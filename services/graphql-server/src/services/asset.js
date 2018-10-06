const nats = require('nats')
const opentracing = require('opentracing')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

class AssetService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('AssetService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'AssetService' })
    this.nats = options.nats || nats.connect({
      servers: config.natsServers,
      user: config.natsUser,
      pass: config.natsPassword
    })
  }

  async exec (context, name, func) {
    let result, err
    let latency

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.SPAN_KIND, 'requester')
    span.setTag(opentracing.Tags.PEER_SERVICE, 'asset-service')
    span.setTag(opentracing.Tags.PEER_ADDRESS, this.nats.currentServer.url.host)

    const carrier = {}
    this.tracer.inject(span, opentracing.FORMAT_TEXT_MAP, carrier)
    const spanContext = JSON.stringify(carrier)

    // Core functionality
    try {
      const startTime = Date.now()
      result = await func(spanContext)
      latency = (Date.now() - startTime) / 1000
    } catch (e) {
      err = e
      this.logger.error(err)
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

  getAsset (context, id) {
    return this.exec(context, 'get-asset', spanContext => {
      return Promise.resolve({ id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' })
    })
  }

  getAssets (context, siteId) {
    return this.exec(context, 'get-assets', spanContext => {
      return Promise.resolve([
        { id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' },
        { id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 }
      ])
    })
  }

  createAlarm (context, input) {
    this.logger.debug(`Received ${input}`)

    return this.exec(context, 'create-alarm', spanContext => {
      return Promise.resolve({ id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' })
    })
  }

  updateAlarm (context, id, input) {
    return this.exec(context, 'update-alarm', spanContext => {
      return Promise.resolve({ id: 'aaaa-aaaa', siteId: '1111-1111', serialNo: '1001', material: 'smoke' })
    })
  }

  deleteAlarm (context, id) {
    return this.exec(context, 'delete-alarm', spanContext => {
      return Promise.resolve(true)
    })
  }

  createCamera (context, input) {
    this.logger.debug(`Received ${input}`)

    return this.exec(context, 'create-camera', spanContext => {
      return Promise.resolve({ id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 })
    })
  }

  updateCamera (context, id, input) {
    return this.exec(context, 'update-camera', spanContext => {
      return Promise.resolve({ id: 'bbbb-bbbb', siteId: '1111-1111', serialNo: '2001', resolution: 921600 })
    })
  }

  deleteCamera (context, id) {
    return this.exec(context, 'delete-camera', spanContext => {
      return Promise.resolve(true)
    })
  }
}

module.exports = AssetService
