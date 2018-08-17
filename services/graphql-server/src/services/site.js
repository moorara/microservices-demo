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
    const site = Object.assign({}, input)
    site.id = _.uniqueId()
    this.store.sites.push(site)
    return Promise.resolve(site)
  }

  all (context, query) {
    const sites = Object.assign([], this.store.sites)
    return Promise.resolve(sites)
  }

  get (context, id) {
    const site = Object.assign({}, this.store.sites.find(s => s.id === id))
    return Promise.resolve(site)
  }

  update (context, id, input) {
    const site = Object.assign({}, { id }, input)
    for (let i in this.store.sites) {
      if (this.store.sites[i].id === id) {
        this.store.sites[i] = site
        const updated = Object.assign({}, site)
        return Promise.resolve(updated)
      }
    }
  }

  modify (context, id, input) {
    const site = Object.assign({}, { id }, input)
    for (let i in this.store.sites) {
      if (this.store.sites[i].id === id) {
        Object.assign(this.store.sites[i], site)
        const updated = Object.assign({}, this.store.sites[i])
        return Promise.resolve(updated)
      }
    }
  }

  delete (context, id) {
    _.remove(this.store.sites, s => s.id === id)
    return Promise.resolve()
  }
}

module.exports = SiteService
