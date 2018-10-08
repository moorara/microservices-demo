/* eslint-env mocha */
const should = require('should')
const promClient = require('prom-client')
const opentracing = require('opentracing')

const GraphQLRouter = require('./router')

describe('GraphQLRouter', () => {
  describe('constructor', () => {
    let config, options

    beforeEach(() => {
      config = {
        graphiQlEnabled: false
      }

      options = {
        logger: {
          trace () {},
          debug () {},
          info () {},
          warn () {},
          error () {},
          fatal () {}
        },
        register: new promClient.Registry(),
        tracer: new opentracing.MockTracer(),
        context: {},
        siteService: {},
        sensorService: {},
        switchService: {},
        assetService: {}
      }
    })

    it('should create a new graphql router with defaults', () => {
      const router = new GraphQLRouter(config, { tracer: options.tracer, switchService: {}, assetService: {} })
      should.exist(router.router)
    })
    it('should create a new graphql router with provided options', () => {
      options.typeDefs = `type Query { hello: String }`
      options.resolvers = {
        Query: {
          hello: () => 'hello, world!'
        }
      }

      const router = new GraphQLRouter(config, options)
      should.exist(router.router)
    })
  })
})
