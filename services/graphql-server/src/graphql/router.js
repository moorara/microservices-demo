const fs = require('fs')
const path = require('path')
const express = require('express')
// const { buildSchema } = require('graphql')
const graphqlHTTP = require('express-graphql')
const { makeExecutableSchema } = require('graphql-tools')

const resolvers = require('./resolvers')
const Logger = require('../utils/logger')
const SiteService = require('../services/site')
const SensorService = require('../services/sensor')
const SwitchService = require('../services/switch')

class GraphQLRouter {
  constructor (config, options) {
    options = options || {}
    const graphiql = config.graphiQlEnabled || false

    const context = Object.assign({}, options.context, {
      logger: options.logger || new Logger('GraphQL'),
      siteService: options.siteService || new SiteService(config, options),
      sensorService: options.sensorService || new SensorService(config, options),
      switchService: options.switchService || new SwitchService(config, options)
    })

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
  }
}

module.exports = GraphQLRouter
