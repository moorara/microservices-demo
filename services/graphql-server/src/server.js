const http = require('http')
const express = require('express')

const ConfigProvider = require('./config')
const Logger = require('./utils/logger')
const Metrics = require('./utils/metrics')
const { createTracer } = require('./utils/tracer')
const middleware = require('./middleware')
const MonitorMiddleware = require('./middleware/monitor')
const livenessRouter = require('./routes/liveness')
const readinessRouter = require('./routes/readiness')
const GraphQLRouter = require('./graphql/router')

class Server {
  constructor (options) {
    options = options || {}

    Logger.addMetadata({
      service: process.env.SERVICE_NAME || 'graphql-server'
    })

    this.configProvider = options.configProvider || new ConfigProvider()
    this.logger = options.logger || new Logger('Server')
    this.metrics = options.metrics || new Metrics()
    this.tracer = options.tracer

    this.app = options.app || express()
    this.routers = options.routers || {}
    this.middleware = options.middleware || {}
  }

  async start () {
    try {
      const config = await this.configProvider.getConfig()

      // Dependencies
      this.tracer = this.tracer || createTracer(config)
      this.routers.liveness = this.routers.liveness || livenessRouter
      this.routers.readiness = this.routers.readiness || readinessRouter

      this.routers.graphql = this.routers.graphql || new GraphQLRouter(config, {
        register: this.metrics.register,
        tracer: this.tracer
      })

      this.middleware.monitor = this.middleware.monitor || new MonitorMiddleware({
        register: this.metrics.register,
        tracer: this.tracer,
        spanName: 'graphql-request'
      })

      // General routes
      this.app.use(this.routers.liveness)
      this.app.use(this.routers.readiness)
      this.app.use(this.metrics.router)

      // Middleware
      this.app.use(this.middleware.monitor.router)

      // Service routes
      this.app.use(this.routers.graphql.router)

      // Make sure unexpected errors are not sent
      this.app.use(middleware.catchError)

      // Create and start a http server
      const server = http.createServer(this.app)
      server.listen(config.servicePort, () => {
        this.logger.info(`🚀 Listening on port ${config.servicePort} ...`)
      })
    } catch (err) {
      this.logger.fatal('Server errored:', err)
      throw err
    }
  }
}

if (process.env.NODE_ENV !== 'test') {
  const server = new Server()
  server.start()
}

module.exports = Server
