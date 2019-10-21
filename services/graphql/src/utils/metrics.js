const express = require('express')
const promClient = require('prom-client')

class Metrics {
  constructor (options) {
    options = options || {}
    this.register = options.register || promClient.register

    this.interval = promClient.collectDefaultMetrics({
      register: this.register
    })

    this.router = express.Router()
    this.router.get('/metrics', this.getMetrics.bind(this))
  }

  getMetrics (req, res) {
    res.type('text')
    res.send(this.register.metrics())
  }
}

module.exports = Metrics
