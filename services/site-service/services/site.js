// http://mongoosejs.com/docs/api.html

const _ = require('lodash')
const opentracing = require('opentracing')

const Logger = require('../util/logger')
const Site = require('../models/site')

class SiteService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SiteService')
    this.tracer = options.tracer || opentracing.globalTracer()
    this.SiteModel = options.SiteModel || new Site().Model
  }

  async create (specs, context) {
    let site

    const siteModel = new this.SiteModel(
      _.pick(specs, [ 'name', 'location', 'tags', 'priority' ])
    )

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('save-document', { childOf: parentSpan })

    try {
      site = await siteModel.save()
      site = site.toJSON()
    } catch (err) {
      this.logger.error(`Error creating site ${specs.name}.`, err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'save-document',
      'db.statement': 'save'
    })
    span.finish()

    return site
  }

  /**
   * query: name, location, tags, minPriority, maxPriority, limit, skip
   */
  async all (query, context) {
    let sites
    let mongoQuery = {}

    query = query || {}
    if (query.name) mongoQuery.name = new RegExp(`.*${query.name}.*`, 'i')
    if (query.location) mongoQuery.location = new RegExp(`.*${query.location}.*`, 'i')
    if (query.tags) _.set(mongoQuery, 'tags.$in', query.tags.split(','))
    if (query.minPriority) _.set(mongoQuery, 'priority.$gte', +query.minPriority)
    if (query.maxPriority) _.set(mongoQuery, 'priority.$lte', +query.maxPriority)

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('find-documents', { childOf: parentSpan })

    try {
      sites = await this.SiteModel
        .find(mongoQuery)
        .limit(+query.limit)
        .skip(+query.skip)
        .exec()
    } catch (err) {
      this.logger.error('Error getting sites.', err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'find-documents',
      'db.statement': 'find'
    })
    span.finish()

    return sites.map(l => l.toJSON())
  }

  async get (id, context) {
    let site

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('find-document', { childOf: parentSpan })

    try {
      site = await this.SiteModel.findById(id)
      site = site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error getting site ${id}.`, err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'find-document',
      'db.statement': 'findById'
    })
    span.finish()

    return site
  }

  async update (id, specs, context) {
    let result

    const query = { _id: id }
    const update = Object.assign({}, _.pick(specs, [ 'name', 'location', 'tags', 'priority' ]))
    const options = { upsert: false, runValidators: true, overwrite: true }

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('update-document', { childOf: parentSpan })

    try {
      let res = await this.SiteModel.update(query, update, options)
      result = res.ok === 1 && res.n === 1
    } catch (err) {
      this.logger.error(`Error updating site ${id}.`, err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'update-document',
      'db.statement': 'update'
    })
    span.finish()

    return result
  }

  async modify (id, specs, context) {
    let site

    const update = Object.assign({}, _.pick(specs, [ 'name', 'location', 'tags', 'priority' ]))
    const options = { new: true, upsert: false, runValidators: true }

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('modify-document', { childOf: parentSpan })

    try {
      site = await this.SiteModel.findByIdAndUpdate(id, update, options)
      site = site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error updating site ${id}.`, err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'modify-document',
      'db.statement': 'findByIdAndUpdate'
    })
    span.finish()

    return site
  }

  async delete (id, context) {
    let site

    context = context || {}
    const parentSpan = context.span
    const span = this.tracer.startSpan('delete-document', { childOf: parentSpan })

    try {
      site = await this.SiteModel.findByIdAndRemove(id)
      site = site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error deleting site ${id}.`, err)
      throw err
    }

    span.setTag('db.type', 'mongo')
    span.log({
      'event': 'delete-document',
      'db.statement': 'findByIdAndRemove'
    })
    span.finish()

    return site
  }
}

module.exports = SiteService
