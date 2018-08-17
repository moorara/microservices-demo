/* eslint-env mocha */
const should = require('should')

const GraphQLRouter = require('./router')

describe('GraphQLRouter', () => {
  describe('constructor', () => {
    let config, logger

    beforeEach(() => {
      config = {
        graphiQlEnabled: false
      }
      logger = {
        trace () {},
        debug () {},
        info () {},
        warn () {},
        error () {},
        fatal () {}
      }
    })

    it('should create a new graphql router with defaults', () => {
      const router = new GraphQLRouter(config)
      should.exist(router.router)
    })
    it('should create a new graphql router with provided options', () => {
      const options = {
        logger,
        context: {},
        siteService: {},
        sensorService: {},
        switchService: {},
        typeDefs: `type Query { hello: String }`,
        resolvers: {
          Query: {
            hello: () => 'hello, world!'
          }
        }
      }

      const router = new GraphQLRouter(config, options)
      should.exist(router.router)
    })
  })
})
