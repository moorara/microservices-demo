const grpc = require('grpc')
const opentracing = require('opentracing')

const { proto } = require('../grpc')
const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

class SwitchService {
  constructor (config, options) {
    options = options || {}
    this.serviceAddr = config.switchServiceAddr
    this.logger = options.logger || new Logger('SwitchService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'SwitchService' })
    this.client = options.client || new proto.SwitchService(
      this.serviceAddr,
      grpc.credentials.createInsecure()
    )
  }

  async exec (context, name, func) {
    let result, err
    let latency

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.SPAN_KIND, 'client')
    span.setTag(opentracing.Tags.PEER_SERVICE, 'switch-service')
    span.setTag(opentracing.Tags.PEER_ADDRESS, this.serviceAddr)

    const carrier = {}
    this.tracer.inject(span, opentracing.FORMAT_TEXT_MAP, carrier)

    const metadata = new grpc.Metadata()
    metadata.add('span.context', JSON.stringify(carrier))

    // Core functionality
    try {
      const startTime = Date.now()
      result = await func(metadata)
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

  installSwitch (context, { siteId, name, state, states }) {
    return this.exec(context, 'install-switch', metadata => {
      return new Promise((resolve, reject) => {
        this.client.installSwitch({ siteId, name, state, states }, metadata, (err, swtch) => {
          return err ? reject(err) : resolve(swtch)
        })
      })
    })
  }

  removeSwitch (context, id) {
    return this.exec(context, 'remove-switch', metadata => {
      return new Promise((resolve, reject) => {
        this.client.removeSwitch({ id }, metadata, err => {
          return err ? reject(err) : resolve()
        })
      })
    })
  }

  getSwitch (context, id) {
    return this.exec(context, 'get-switch', metadata => {
      return new Promise((resolve, reject) => {
        this.client.getSwitch({ id }, metadata, (err, swtch) => {
          return err ? reject(err) : resolve(swtch)
        })
      })
    })
  }

  getSwitches (context, siteId) {
    return this.exec(context, 'get-switches', metadata => {
      return new Promise((resolve, reject) => {
        const switches = []
        const call = this.client.getSwitches({ siteId }, metadata)
        call.on('data', swtch => switches.push(swtch))
        call.on('error', err => reject(err))
        call.on('end', () => resolve(switches))
      })
    })
  }

  setSwitch (context, id, { state }) {
    return this.exec(context, 'set-switch', metadata => {
      return new Promise((resolve, reject) => {
        this.client.setSwitch({ id, state }, metadata, err => {
          if (err) {
            reject(err)
          }

          // switch-service does not respond with updated switch
          this.client.getSwitch({ id }, metadata, (err, swtch) => {
            return err ? reject(err) : resolve(swtch)
          })
        })
      })
    })
  }
}

module.exports = SwitchService
