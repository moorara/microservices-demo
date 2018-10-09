const nats = require('nats')
const opentracing = require('opentracing')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

const subject = 'asset_service'
const timeout = 1000

class AssetService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('AssetService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'AssetService' })
    this.conn = options.conn || nats.connect({
      servers: config.natsServers,
      user: config.natsUser,
      pass: config.natsPassword
    })
  }

  async exec (context, name, func) {
    let err, result
    let startTime, latency

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.SPAN_KIND, 'requester')
    span.setTag(opentracing.Tags.PEER_SERVICE, 'asset-service')
    span.setTag(opentracing.Tags.PEER_ADDRESS, this.conn.currentServer.url.host)

    const carrier = {}
    this.tracer.inject(span, opentracing.FORMAT_TEXT_MAP, carrier)
    const spanContext = JSON.stringify(carrier)

    // Core functionality
    try {
      startTime = Date.now()
      result = await func(spanContext)
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

  sendRequest (request, callback) {
    this.conn.requestOne(subject, JSON.stringify(request), {}, timeout, resp => {
      if (resp instanceof nats.NatsError && resp.code === nats.REQ_TIMEOUT) {
        return callback(resp)
      }

      let response

      try {
        response = JSON.parse(resp)
      } catch (err) {
        return callback(err)
      }

      if (response.kind !== request.kind) {
        return callback(new Error(`Invalid response ${response.kind}`))
      }

      callback(null, response)
    })
  }

  createAlarm (context, input) {
    return this.exec(context, 'create-alarm', span => {
      return new Promise((resolve, reject) => {
        const kind = 'createAlarm'
        const request = { kind, span, input }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.alarm)
        })
      })
    })
  }

  allAlarm (context, siteId) {
    return this.exec(context, 'all-alarm', span => {
      return new Promise((resolve, reject) => {
        const kind = 'allAlarm'
        const request = { kind, span, siteId }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.alarms)
        })
      })
    })
  }

  getAlarm (context, id) {
    return this.exec(context, 'get-alarm', span => {
      return new Promise((resolve, reject) => {
        const kind = 'getAlarm'
        const request = { kind, span, id }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.alarm)
        })
      })
    })
  }

  updateAlarm (context, id, input) {
    return this.exec(context, 'update-alarm', span => {
      return new Promise((resolve, reject) => {
        const kind = 'updateAlarm'
        const request = { kind, span, id, input }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.updated)
        })
      })
    })
  }

  deleteAlarm (context, id) {
    return this.exec(context, 'delete-alarm', span => {
      return new Promise((resolve, reject) => {
        const kind = 'deleteAlarm'
        const request = { kind, span, id }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.deleted)
        })
      })
    })
  }

  createCamera (context, input) {
    return this.exec(context, 'create-camera', span => {
      return new Promise((resolve, reject) => {
        const kind = 'createCamera'
        const request = { kind, span, input }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.camera)
        })
      })
    })
  }

  allCamera (context, siteId) {
    return this.exec(context, 'all-camera', span => {
      return new Promise((resolve, reject) => {
        const kind = 'allCamera'
        const request = { kind, span, siteId }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.cameras)
        })
      })
    })
  }

  getCamera (context, id) {
    return this.exec(context, 'get-camera', span => {
      return new Promise((resolve, reject) => {
        const kind = 'getCamera'
        const request = { kind, span, id }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.camera)
        })
      })
    })
  }

  updateCamera (context, id, input) {
    return this.exec(context, 'update-camera', span => {
      return new Promise((resolve, reject) => {
        const kind = 'updateCamera'
        const request = { kind, span, id, input }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.updated)
        })
      })
    })
  }

  deleteCamera (context, id) {
    return this.exec(context, 'delete-camera', span => {
      return new Promise((resolve, reject) => {
        const kind = 'deleteCamera'
        const request = { kind, span, id }

        this.sendRequest(request, (err, response) => {
          if (err) {
            return reject(err)
          }
          resolve(response.deleted)
        })
      })
    })
  }
}

module.exports = AssetService
