const express = require('express')
const bodyParser = require('body-parser')

const Logger = require('../util/logger')
const Middleware = require('../middleware')
const LinkService = require('../services/link')

class LinkRouter {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('LinkRouter')
    this.linkService = options.linkService || new LinkService(config, options)

    this.router = express.Router()
    this.router.use(bodyParser.json())

    this.router.post('*', Middleware.ensureJson())
    this.router.put('*', Middleware.ensureJson())

    this.router.route('/')
      .post(this.postLink.bind(this))
      .get(this.getLinks.bind(this))

    this.router.route('/:id')
      .get(this.getLink.bind(this))
      .put(this.putLink.bind(this))
      .delete(this.deleteLink.bind(this))
  }

  async postLink (req, res, next) {
    let link
    let specs = req.body

    try {
      link = await this.linkService.create(specs)
      res.status(201).send(link)
    } catch (err) {
      this.logger.error('Failed to create new link.', err)
      return next(err)
    }
  }

  async getLinks (req, res, next) {
    let links
    let query = req.query

    try {
      links = await this.linkService.getAll(query)
      res.status(200).send(links)
    } catch (err) {
      this.logger.error('Failed to get links.', err)
      return next(err)
    }
  }

  async getLink (req, res, next) {
    let link
    let id = req.params.id

    try {
      link = await this.linkService.get(id)
    } catch (err) {
      this.logger.error(`Failed to get link ${id}.`, err)
      return next(err)
    }

    if (!link) {
      res.sendStatus(404)
    } else {
      res.status(200).send(link)
    }
  }

  async putLink (req, res, next) {
    let link
    let id = req.params.id
    let specs = req.body

    try {
      link = await this.linkService.update(id, specs)
    } catch (err) {
      this.logger.error(`Failed to update link ${id}.`, err)
      return next(err)
    }

    if (!link) {
      res.sendStatus(404)
    } else {
      res.status(200).send(link)
    }
  }

  async deleteLink (req, res, next) {
    let link
    let id = req.params.id

    try {
      link = await this.linkService.delete(id)
    } catch (err) {
      this.logger.error(`Failed to delete link ${id}.`, err)
      return next(err)
    }

    if (!link) {
      res.sendStatus(404)
    } else {
      res.status(200).send(link)
    }
  }
}

module.exports = LinkRouter
