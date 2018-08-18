const _ = require('lodash')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

class SensorService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SensorService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'sensor-service' })

    this.store = {
      sensors: [
        { id: '1111-1111', siteId: 'aaaa-aaaa', name: 'temperature', unit: 'celsius', minSafe: -30.0, maxSafe: 30.0 },
        { id: '2222-2222', siteId: 'bbbb-bbbb', name: 'temperature', unit: 'fahrenheit', minSafe: -22.0, maxSafe: 86.0 }
      ]
    }
  }

  create (context, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('create-sensor', { childOf: context.span })

    // TODO
    const sensor = Object.assign({}, input)
    sensor.id = _.uniqueId()
    this.store.sensors.push(sensor)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'create_sensor', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'sensor-service')
    span.setTag('peer.address', 'sensor-service:4020')
    span.log({ event: 'create-sensor', message: '' })
    span.finish()

    return Promise.resolve(sensor)
  }

  all (context, siteId) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-sensors', { childOf: context.span })

    // TODO
    const sensors = this.store.sensors.filter(s => s.siteId === siteId)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_sensors', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'sensor-service')
    span.setTag('peer.address', 'sensor-service:4020')
    span.log({ event: 'get-sensors', message: '' })
    span.finish()

    return Promise.resolve(sensors)
  }

  get (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-sensor', { childOf: context.span })

    // TODO
    const sensor = Object.assign({}, this.store.sensors.find(s => s.id === id))

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_sensor', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'sensor-service')
    span.setTag('peer.address', 'sensor-service:4020')
    span.log({ event: 'get-sensor', message: '' })
    span.finish()

    return Promise.resolve(sensor)
  }

  update (context, id, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('update-sensor', { childOf: context.span })

    // TODO
    let updated
    const sensor = Object.assign({}, { id }, input)
    for (let i in this.store.sensors) {
      if (this.store.sensors[i].id === id) {
        this.store.sensors[i] = sensor
        updated = Object.assign({}, sensor)
        break
      }
    }

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'update_sensor', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'sensor-service')
    span.setTag('peer.address', 'sensor-service:4020')
    span.log({ event: 'update-sensor', message: '' })
    span.finish()

    return Promise.resolve(updated)
  }

  delete (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('delete-sensor', { childOf: context.span })

    // TODO
    _.remove(this.store.sensors, s => s.id === id)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'delete_sensor', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'sensor-service')
    span.setTag('peer.address', 'sensor-service:4020')
    span.log({ event: 'delete-sensor', message: '' })
    span.finish()

    return Promise.resolve()
  }
}

module.exports = SensorService
