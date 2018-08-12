const http = require('http')
const express = require('express')

const ConfigProvider = require('./config')
const Logger = require('./utils/logger')
const Metrics = require('./utils/metrics')
const Tracer = require('./utils/tracer')
const LivenessRouter = require('./routes/liveness')
const ReadinessRouter = require('./routes/readiness')
const Middleware = require('./middleware')
const MonitorMiddleware = require('./middleware/monitor')

class Server {
  constructor (options) {
    Logger.addMetadata({
      service: process.env.SERVICE_NAME || 'graphql-server'
    })

    options = options || {}
    this.app = options.app || express()
    this.routers = options.routers || {}
    this.logger = options.logger || new Logger('Server')
    this.metrics = options.metrics || new Metrics()
    this.tracer = options.tracer
    this.configProvider = options.configProvider || new ConfigProvider()
  }

  async start () {
    try {
      const config = await this.configProvider.getConfig()

      // Dependencies
      this.tracer = this.tracer || Tracer.createTracer(config)

      // General routes
      this.app.use(LivenessRouter)
      this.app.use(ReadinessRouter)
      this.app.use(this.metrics.router)

      // Monitoring middleware for logs, metrics, and traces
      const monitorMiddleware = new MonitorMiddleware({
        register: this.metrics.register,
        tracer: this.tracer,
        spanName: 'graphql-request'
      })
      this.app.use(monitorMiddleware.router)

      // Service routes
      this.app.use((req, res, next) => res.sendStatus(200))

      // Make sure unexpected errors are not sent
      this.app.use(Middleware.catchError)

      // Create and start a http server
      const server = http.createServer(this.app)
      server.listen(config.servicePort, () => {
        this.logger.info(`Listening on port ${config.servicePort} ...`)
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
