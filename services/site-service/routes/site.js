const express = require('express')
const bodyParser = require('body-parser')

const Logger = require('../util/logger')
const Middleware = require('../middleware')
const SiteService = require('../services/site')

class SiteRouter {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SiteRouter')
    this.siteService = options.siteService || new SiteService(config, options)

    this.router = express.Router()
    this.router.use(bodyParser.json())

    this.router.post('*', Middleware.ensureJson())
    this.router.put('*', Middleware.ensureJson())

    this.router.route('/')
      .post(this.postSite.bind(this))
      .get(this.getSites.bind(this))

    this.router.route('/:id')
      .get(this.getSite.bind(this))
      .put(this.putSite.bind(this))
      .delete(this.deleteSite.bind(this))
  }

  async postSite (req, res, next) {
    let site
    let specs = req.body

    try {
      site = await this.siteService.create(specs)
      res.status(201).send(site)
    } catch (err) {
      this.logger.error('Failed to create new site.', err)
      return next(err)
    }
  }

  async getSites (req, res, next) {
    let sites
    let query = req.query

    try {
      sites = await this.siteService.all(query)
      res.status(200).send(sites)
    } catch (err) {
      this.logger.error('Failed to get sites.', err)
      return next(err)
    }
  }

  async getSite (req, res, next) {
    let site
    let id = req.params.id

    try {
      site = await this.siteService.get(id)
    } catch (err) {
      this.logger.error(`Failed to get site ${id}.`, err)
      return next(err)
    }

    if (!site) {
      res.sendStatus(404)
    } else {
      res.status(200).send(site)
    }
  }

  async putSite (req, res, next) {
    let site
    let id = req.params.id
    let specs = req.body

    try {
      site = await this.siteService.update(id, specs)
    } catch (err) {
      this.logger.error(`Failed to update site ${id}.`, err)
      return next(err)
    }

    if (!site) {
      res.sendStatus(404)
    } else {
      res.status(200).send(site)
    }
  }

  async deleteSite (req, res, next) {
    let site
    let id = req.params.id

    try {
      site = await this.siteService.delete(id)
    } catch (err) {
      this.logger.error(`Failed to delete site ${id}.`, err)
      return next(err)
    }

    if (!site) {
      res.sendStatus(404)
    } else {
      res.status(200).send(site)
    }
  }
}

module.exports = SiteRouter
