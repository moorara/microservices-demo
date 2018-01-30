const http = require('http')
const express = require('express')

const Logger = require('./util/logger')
const ConfigProvider = require('./config')
const Mongo = require('./models/mongo')
const Middleware = require('./middleware')
const MetricsMiddleware = require('./middleware/metrics')
const LoggerMiddleware = require('./middleware/logger')
const HealthRouter = require('./routes/health')
const LinkRouter = require('./routes/link')

class Server {
  constructor (options) {
    Logger.addContext({
      service: process.env.SERVICE_NAME || 'node-service'
    })

    options = options || {}
    this.app = options.app || express()
    this.routers = options.routers || {}
    this.mongo = options.mongo
    this.logger = options.logger || new Logger('Server')
    this.configProvider = options.configProvider || new ConfigProvider()
  }

  async start () {
    try {
      const config = await this.configProvider.getConfig()

      // Dependencies
      this.mongo = this.mongo || new Mongo(config)
      this.metricsMiddleware = new MetricsMiddleware()
      this.routers.links = this.routers.links || new LinkRouter(config)

      // Unauthenticated routes
      this.app.use('/health', HealthRouter)

      this.app.use(Middleware.normalize())
      this.app.use(this.metricsMiddleware.router)
      this.app.use(LoggerMiddleware.http())

      // Authenticated routes
      this.app.use('/v1/links', this.routers.links.router)

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
