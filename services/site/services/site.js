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

  async exec (context, name, mongoQuery, func) {
    context = context || {}
    let latency, result, err

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.DB_TYPE, 'mongo')
    span.setTag(opentracing.Tags.DB_STATEMENT, mongoQuery)

    // Core functionality
    try {
      const startTime = Date.now()
      result = await func()
      latency = (Date.now() - startTime) / 1000
    } catch (e) {
      err = e
      this.logger.error(err)
    }

    // Metrics

    // Tracing
    span.log({
      event: name,
      latency: latency,
      message: err ? err.message : 'successful!'
    })
    span.finish()

    if (err) {
      throw err
    }

    return result
  }

  async create (specs, context) {
    return this.exec(context, 'save-document', 'save()', async () => {
      const siteModel = new this.SiteModel(
        _.pick(specs, ['name', 'location', 'tags', 'priority'])
      )

      const site = await siteModel.save()
      return site.toJSON()
    })
  }

  /**
   * query: name, location, tags, minPriority, maxPriority, limit, skip
   */
  async all (query, context) {
    return this.exec(context, 'find-documents', 'find().limit().skip().exec()', async () => {
      query = query || {}
      const mongoQuery = {}

      if (query.name) mongoQuery.name = new RegExp(`.*${query.name}.*`, 'i')
      if (query.location) mongoQuery.location = new RegExp(`.*${query.location}.*`, 'i')
      if (query.tags) _.set(mongoQuery, 'tags.$in', query.tags.split(','))
      if (query.minPriority) _.set(mongoQuery, 'priority.$gte', +query.minPriority)
      if (query.maxPriority) _.set(mongoQuery, 'priority.$lte', +query.maxPriority)

      const sites = await this.SiteModel.find(mongoQuery).limit(+query.limit).skip(+query.skip).exec()
      return sites.map(l => l.toJSON())
    })
  }

  async get (id, context) {
    return this.exec(context, 'find-document', 'findById()', async () => {
      const site = await this.SiteModel.findById(id)
      return site ? site.toJSON() : null
    })
  }

  async update (id, specs, context) {
    return this.exec(context, 'update-document', 'update()', async () => {
      const query = { _id: id }
      const update = Object.assign({}, _.pick(specs, ['name', 'location', 'tags', 'priority']))
      const options = { upsert: false, runValidators: true, overwrite: true }

      const res = await this.SiteModel.update(query, update, options)
      return res.ok === 1 && res.n === 1
    })
  }

  async modify (id, specs, context) {
    return this.exec(context, 'modify-document', 'findByIdAndUpdate()', async () => {
      const update = Object.assign({}, _.pick(specs, ['name', 'location', 'tags', 'priority']))
      const options = { new: true, upsert: false, runValidators: true }

      const site = await this.SiteModel.findByIdAndUpdate(id, update, options)
      return site ? site.toJSON() : null
    })
  }

  async delete (id, context) {
    return this.exec(context, 'delete-document', 'findByIdAndRemove()', async () => {
      const site = await this.SiteModel.findByIdAndRemove(id)
      return site ? site.toJSON() : null
    })
  }
}

module.exports = SiteService
