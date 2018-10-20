const http = require('http')
const express = require('express')

const Logger = require('./util/logger')
const Metrics = require('./util/metrics')
const { createTracer } = require('./util/tracer')
const ConfigProvider = require('./config')
const Mongo = require('./models/mongo')
const Middleware = require('./middleware')
const LoggerMiddleware = require('./middleware/logger')
const MetricsMiddleware = require('./middleware/metrics')
const TracerMiddleware = require('./middleware/tracer')
const HealthRouter = require('./routes/health')
const SiteRouter = require('./routes/site')

class Server {
  constructor (options) {
    Logger.addMetadata({
      service: process.env.SERVICE_NAME || 'site-service'
    })

    options = options || {}
    this.app = options.app || express()
    this.routers = options.routers || {}
    this.mongo = options.mongo
    this.logger = options.logger || new Logger('Server')
    this.metrics = options.metrics || new Metrics()
    this.tracer = options.tracer
    this.configProvider = options.configProvider || new ConfigProvider()
  }

  async start () {
    try {
      const config = await this.configProvider.getConfig()

      // Dependencies
      this.mongo = this.mongo || new Mongo(config)
      this.tracer = this.tracer || createTracer(config)
      this.routers.sites = this.routers.sites || new SiteRouter(config, { tracer: this.tracer })

      // Unauthenticated routes
      this.app.use('/health', HealthRouter)
      this.app.use(this.metrics.router)

      this.app.use(Middleware.normalize())
      this.app.use(LoggerMiddleware.create())
      this.app.use(MetricsMiddleware.create({ register: this.metrics.register }))
      this.app.use(TracerMiddleware.create({ tracer: this.tracer }))

      // Authenticated routes
      this.app.use('/v1/sites', this.routers.sites.router)

      this.app.use(Middleware.catchError({
        environment: process.env.NODE_ENV
      }))

      const server = http.createServer(this.app)

      // Connect to Mongo
      const connection = await this.mongo.connect()
      connection.once('error', err => {
        this.logger.error('Mongo connection error.', err)
        server.close()
      })

      server.listen(config.servicePort, () => {
        this.logger.info(`Listening on port ${config.servicePort} ...`)
      })
    } catch (err) {
      this.logger.fatal(err)
      throw err
    }
  }
}

if (process.env.NODE_ENV !== 'test') {
  const server = new Server()
  server.start()
}

module.exports = Server
