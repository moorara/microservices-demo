const _ = require('lodash')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

class SiteService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SiteService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'site-service' })

    this.store = {
      sites: [
        { id: 'aaaa-aaaa', name: 'Gas Station', location: 'Toronto, ON', tags: ['energy', 'gas'] },
        { id: 'bbbb-bbbb', name: 'Power Plant', location: 'Montreal, QC', tags: ['energy', 'power'] }
      ]
    }
  }

  create (context, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('create-site', { childOf: context.span })

    // TODO
    const site = Object.assign({}, input)
    site.id = _.uniqueId()
    this.store.sites.push(site)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'create_site', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'create-site', message: '' })
    span.finish()

    return Promise.resolve(site)
  }

  all (context, query) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-sites', { childOf: context.span })

    // TODO
    const sites = Object.assign([], this.store.sites)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_sites', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'get-sites', message: '' })
    span.finish()

    return Promise.resolve(sites)
  }

  get (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('get-site', { childOf: context.span })

    // TODO
    const site = Object.assign({}, this.store.sites.find(s => s.id === id))

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'get_site', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'get-site', message: '' })
    span.finish()

    return Promise.resolve(site)
  }

  update (context, id, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('update-site', { childOf: context.span })

    // TODO
    let updated
    const site = Object.assign({}, { id }, input)
    for (let i in this.store.sites) {
      if (this.store.sites[i].id === id) {
        this.store.sites[i] = site
        updated = Object.assign({}, site)
        break
      }
    }

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'update_site', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'update-site', message: '' })
    span.finish()

    return Promise.resolve(updated)
  }

  modify (context, id, input) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('modify-site', { childOf: context.span })

    // TODO
    let updated
    const site = Object.assign({}, { id }, input)
    for (let i in this.store.sites) {
      if (this.store.sites[i].id === id) {
        Object.assign(this.store.sites[i], site)
        updated = Object.assign({}, this.store.sites[i])
        break
      }
    }

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'modify_site', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'modify-site', message: '' })
    span.finish()

    return Promise.resolve(updated)
  }

  delete (context, id) {
    const startTime = +new Date()
    const span = this.tracer.startSpan('delete-site', { childOf: context.span })

    // TODO
    _.remove(this.store.sites, s => s.id === id)

    const endTime = +new Date()
    const latency = (endTime - startTime) / 1000

    // Metrics
    const labelValues = { op: 'delete_site', success: 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Traces
    span.setTag('span.kind', 'client')
    span.setTag('peer.service', 'site-service')
    span.setTag('peer.address', 'site-service:4010')
    span.log({ event: 'delete-site', message: '' })
    span.finish()

    return Promise.resolve()
  }
}

module.exports = SiteService
