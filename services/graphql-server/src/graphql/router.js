const fs = require('fs')
const path = require('path')
const express = require('express')
const promClient = require('prom-client')
// const { buildSchema } = require('graphql')
const graphqlHTTP = require('express-graphql')
const { makeExecutableSchema } = require('graphql-tools')
const { default: expressPlayground } = require('graphql-playground-middleware-express')

const resolvers = require('./resolvers')
const Logger = require('../utils/logger')
const SiteService = require('../services/site')
const SensorService = require('../services/sensor')
const SwitchService = require('../services/switch')
const AssetService = require('../services/asset')

const histogramName = 'graphql_operations_latency_seconds'
const histogramHelp = 'latency histogram of graphql operations'
const summaryName = 'graphql_operations_latency_quantiles_seconds'
const summaryHelp = 'latency summary of graphql operations'
const labelNames = [ 'op', 'success' ]
const buckets = [ 0.01, 0.1, 0.5, 1 ]
const percentiles = [ 0.1, 0.5, 0.95, 0.99 ]

class GraphQLRouter {
  /**
   * @param {object}  config                  Configuration
   * @param {boolean} config.graphiQlEnabled  Enable GraphiQL UI?
   * @param {object}  options                 Optional
   * @param {object}  options.logger          A Logger instance
   * @param {object}  options.register        A Prometheus registry
   * @param {object}  options.tracer          A Tracer instance
   * @param {object}  options.context         The context for resolver functions
   * @param {object}  options.siteService     An instance of SiteService
   * @param {object}  options.sensorService   An instance of SensorService
   * @param {object}  options.switchService   An instance of SwitchService
   * @param {object}  options.assetService    An instance of AssetService
   */
  constructor (config, options) {
    options = options || {}
    const graphiql = config.graphiQlEnabled || false
    const register = options.register || promClient.register

    // Histogram metric for downstream services
    options.histogram = new promClient.Histogram({
      name: histogramName,
      help: histogramHelp,
      labelNames,
      buckets,
      registers: [ register ]
    })

    // Summary metric for downstream services
    options.summary = new promClient.Summary({
      name: summaryName,
      help: summaryHelp,
      labelNames,
      percentiles,
      registers: [ register ]
    })

    // Context for resolver functions
    const context = Object.assign({}, options.context, {
      logger: options.logger || new Logger('Resolvers'),
      siteService: options.siteService || new SiteService(config, options),
      sensorService: options.sensorService || new SensorService(config, options),
      switchService: options.switchService || new SwitchService(config, options),
      assetService: options.assetService || new AssetService(config, options)
    })

    // GraphQL schema
    const typeDefs = options.typeDefs || fs.readFileSync(path.join(__dirname, 'schema.graphql'), 'utf-8')
    const schema = makeExecutableSchema({
      typeDefs,
      resolvers: options.resolvers || resolvers
    })

    // Set up router
    this.router = express.Router()
    this.router.use('/graphql', graphqlHTTP(async (req, res, graphQlParams) => {
      Object.assign(context, req.context)
      return { schema, context, graphiql }
    }))
    this.router.use('/playground', expressPlayground({
      endpoint: '/graphql',
      settings: {
        'editor.fontSize': 12,
        'editor.theme': 'dark' // light
      }
    }))
  }
}

module.exports = GraphQLRouter
