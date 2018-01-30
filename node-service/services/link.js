const _ = require('lodash')

const Logger = require('../util/logger')
const Link = require('../models/link')

class LinkService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('LinkService')
    this.LinkModel = options.LinkModel || new Link().Model
  }

  async create (specs) {
    const link = new this.LinkModel(
      _.pick(specs, [ 'url', 'title', 'tags', 'rank' ])
    )

    try {
      let savedLink = await link.save()
      return savedLink.toJSON()
    } catch (err) {
      this.logger.error(`Error creating link ${specs.url}.`, err)
      throw err
    }
  }

  /**
   * query: url, title, tags, minRank, maxRank, limit, skip
   */
  async getAll (query) {
    let links
    let mongoQuery = {}

    query = query || {}
    if (query.url) mongoQuery.url = new RegExp(`.*${query.url}.*`, 'i')
    if (query.title) mongoQuery.title = new RegExp(`.*${query.title}.*`, 'i')
    if (query.tags) _.set(mongoQuery, 'tags.$in', query.tags.split(','))
    if (query.minRank) _.set(mongoQuery, 'rank.$gte', +query.minRank)
    if (query.maxRank) _.set(mongoQuery, 'rank.$lte', +query.maxRank)

    try {
      links = await this.LinkModel
        .find(mongoQuery)
        .limit(+query.limit)
        .skip(+query.skip)
        .exec()
    } catch (err) {
      this.logger.error('Error getting links.', err)
      throw err
    }

    return links.map(l => l.toJSON())
  }

  async get (id) {
    try {
      let link = await this.LinkModel.findById(id)
      return link ? link.toJSON() : null
    } catch (err) {
      this.logger.error(`Error getting link ${id}.`, err)
      throw err
    }
  }

  async update (id, specs) {
    let query = {}
    query = _.merge(query, _.pick(specs, [ 'url', 'title', 'tags', 'rank' ]))

    try {
      let link = await this.LinkModel.findByIdAndUpdate(id, query, { new: true })
      return link ? link.toJSON() : null
    } catch (err) {
      this.logger.error(`Error updating link ${id}.`, err)
      throw err
    }
  }

  async delete (id) {
    try {
      let link = await this.LinkModel.findByIdAndRemove(id)
      return link ? link.toJSON() : null
    } catch (err) {
      this.logger.error(`Error deleting link ${id}.`, err)
      throw err
    }
  }
}

module.exports = LinkService
