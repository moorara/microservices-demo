const _ = require('lodash')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

class SwitchService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SwitchService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'switch-service' })

    this.store = {
      switches: [
        { id: '3333-3333', siteId: 'aaaa-aaaa', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] },
        { id: '4444-4444', siteId: 'bbbb-bbbb', name: 'Light', state: 'OFF', states: ['OFF', 'ON'] }
      ]
    }
  }

  create (context, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('create-switch', { childOf: context.span })

    // TODO
    const swtch = Object.assign({}, input)
    swtch.id = _.uniqueId()
    this.store.switches.push(swtch)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'create_switch', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'switch-service')
    span.setTag('peer.address', 'switch-service:4030')
    span.log({ event: 'create-switch', message: '' })
    span.finish()

    return Promise.resolve(swtch)
  }

  all (context, siteId) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-switches', { childOf: context.span })

    // TODO
    const switches = this.store.switches.filter(s => s.siteId === siteId)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_switches', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'switch-service')
    span.setTag('peer.address', 'switch-service:4030')
    span.log({ event: 'get-switches', message: '' })
    span.finish()

    return Promise.resolve(switches)
  }

  get (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-switch', { childOf: context.span })

    // TODO
    const swtch = Object.assign({}, this.store.switches.find(s => s.id === id))

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_switch', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'switch-service')
    span.setTag('peer.address', 'switch-service:4030')
    span.log({ event: 'get-switch', message: '' })
    span.finish()

    return Promise.resolve(swtch)
  }

  update (context, id, { state }) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('update-switch', { childOf: context.span })

    // TODO
    const swtch = this.store.switches.find(s => s.id === id)
    swtch.state = state
    const updated = Object.assign({}, swtch)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'update_switch', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'switch-service')
    span.setTag('peer.address', 'switch-service:4030')
    span.log({ event: 'update-switch', message: '' })
    span.finish()

    return Promise.resolve(updated)
  }

  delete (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('delete-switch', { childOf: context.span })

    // TODO
    _.remove(this.store.switches, s => s.id === id)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'delete_switch', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'switch-service')
    span.setTag('peer.address', 'switch-service:4030')
    span.log({ event: 'delete-switch', message: '' })
    span.finish()

    return Promise.resolve()
  }
}

module.exports = SwitchService
