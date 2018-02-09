const _ = require('lodash')

const Logger = require('../util/logger')
const Site = require('../models/site')

class SiteService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SiteService')
    this.SiteModel = options.SiteModel || new Site().Model
  }

  async create (specs) {
    const site = new this.SiteModel(
      _.pick(specs, [ 'name', 'location', 'tags', 'priority' ])
    )

    try {
      let savedSite = await site.save()
      return savedSite.toJSON()
    } catch (err) {
      this.logger.error(`Error creating site ${specs.name}.`, err)
      throw err
    }
  }

  /**
   * query: name, location, tags, minPriority, maxPriority, limit, skip
   */
  async all (query) {
    let sites
    let mongoQuery = {}

    query = query || {}
    if (query.name) mongoQuery.name = new RegExp(`.*${query.name}.*`, 'i')
    if (query.location) mongoQuery.location = new RegExp(`.*${query.location}.*`, 'i')
    if (query.tags) _.set(mongoQuery, 'tags.$in', query.tags.split(','))
    if (query.minPriority) _.set(mongoQuery, 'priority.$gte', +query.minPriority)
    if (query.maxPriority) _.set(mongoQuery, 'priority.$lte', +query.maxPriority)

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

    return sites.map(l => l.toJSON())
  }

  async get (id) {
    try {
      let site = await this.SiteModel.findById(id)
      return site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error getting site ${id}.`, err)
      throw err
    }
  }

  async update (id, specs) {
    let query = {}
    query = _.merge(query, _.pick(specs, [ 'name', 'location', 'tags', 'priority' ]))

    try {
      let site = await this.SiteModel.findByIdAndUpdate(id, query, { new: true })
      return site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error updating site ${id}.`, err)
      throw err
    }
  }

  async delete (id) {
    try {
      let site = await this.SiteModel.findByIdAndRemove(id)
      return site ? site.toJSON() : null
    } catch (err) {
      this.logger.error(`Error deleting site ${id}.`, err)
      throw err
    }
  }
}

module.exports = SiteService
