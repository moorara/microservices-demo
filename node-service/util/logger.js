const _ = require('lodash')
const winston = require('winston')

const LOG_LEVELS = { fatal: 0, error: 1, warn: 2, info: 3, debug: 4, trace: 5 }

/**
 * options:
 *   - context: { ... }
 *   - winston: a winston logger instance
 */
class Logger {
  constructor (logger, options) {
    options = options || {}
    this.winston = options.winston || Logger.getWinstonLogger()
    this.context = Object.assign({ pid: process.pid }, Logger.context, options.context, logger ? { logger } : {})
  }

  _log (level, args) {
    let values = []
    let context = Object.assign({}, this.context)

    // Make sure all objects are logged as context
    for (let arg of args) {
      if (!_.isObject(arg)) {
        values.push(arg)
      } else if (arg instanceof Error) {
        values.push(arg.message)
        Object.assign(context, { error: arg })
      } else {
        Object.assign(context, arg)
      }
    }

    this.winston.log(level, ...values, context)
  }

  trace () {
    this._log('trace', [].slice.call(arguments))
  }

  debug () {
    this._log('debug', [].slice.call(arguments))
  }

  info () {
    this._log('info', [].slice.call(arguments))
  }

  warn () {
    this._log('warn', [].slice.call(arguments))
  }

  error () {
    this._log('error', [].slice.call(arguments))
  }

  fatal () {
    this._log('fatal', [].slice.call(arguments))
  }
}

Logger.context = {}

Logger.addContext = function (context) {
  Object.assign(Logger.context, context)
}

Logger.getWinstonLogger = function () {
  if (Logger.winston) {
    return Logger.winston
  }

  const level = process.env.LOG_LEVEL || 'info'

  let opts = {
    level,
    json: true,
    stringify: true,
    handleExceptions: true,
    humanReadableUnhandledException: true,
    timestamp: () => new Date().toISOString()
  }

  if (process.env.NODE_ENV === 'development') {
    opts.json = false
    opts.stringify = false
    opts.colorize = true
    opts.prettyPrint = true
  }

  const transports = [
    process.env.NODE_ENV !== 'test'
      ? new winston.transports.Console(opts)
      : new winston.transports.File({ filename: '/dev/null' })
  ]

  Logger.winston = new winston.Logger({
    transports,
    levels: LOG_LEVELS
  })

  return Logger.winston
}

module.exports = Logger
