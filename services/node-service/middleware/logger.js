const expressWinston = require('express-winston')

const Logger = require('../util/logger')

const defaultSkip = (req, res) => false
const defaultIgnoreRoute = (req, res) => false
const defaultIgnoredRoutes = [ '/health', '/metrics' ]

const defaultRequestFilter = (req, propName) => req[propName]
const defaultRequestWhitelist = [
  'method', 'endpoint', 'params', 'query',
  'httpVersion', 'headers.user-agent'
]

const defaultResponseFilter = (res, propName) => res[propName]
const defaultResponseWhitelist = [ 'statusCode', 'statusClass' ]

module.exports = {
  http (options) {
    options = options || {}
    options.winston = options.winston || Logger.getWinstonLogger()
    let context = Object.assign({}, Logger.context, { logger: 'HttpMiddleware' })

    let loggerMiddleware = expressWinston.logger({
      winstonInstance: options.winston,
      expressFormat: process.env.NODE_ENV === 'development',

      statusLevels: true,
      meta: true,
      baseMeta: context,

      skip: options.skip || defaultSkip,
      ignoreRoute: options.ignoreRoute || defaultIgnoreRoute,
      ignoredRoutes: options.ignoredRoutes || defaultIgnoredRoutes,

      requestFilter: options.requestFilter || defaultRequestFilter,
      requestWhitelist: options.requestWhitelist || defaultRequestWhitelist,

      responseFilter: options.responseFilter || defaultResponseFilter,
      responseWhitelist: options.responseWhitelist || defaultResponseWhitelist
    })

    return loggerMiddleware
  }
}
