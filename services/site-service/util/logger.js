const _ = require('lodash')
const winston = require('winston')

const levels = { fatal: 0, error: 1, warn: 2, info: 3, debug: 4, trace: 5 }

class Logger {
  static addContext (context) {
    Object.assign(this.context, context)
  }

  // ref: https://github.com/winstonjs/winston/tree/2.x
  static getWinstonLogger () {
    if (this.winston) {
      return this.winston
    }

    const level = process.env.LOG_LEVEL || 'info'

    let opts = process.env.NODE_ENV === 'development' ? {
      json: false,
      stringify: false,
      colorize: true,
      prettyPrint: true
    } : {
      json: true,
      stringify: true,
      handleExceptions: true,
      humanReadableUnhandledException: true,
      timestamp: () => new Date().toISOString()
    }

    const transports = [
      process.env.NODE_ENV !== 'test'
        ? new winston.transports.Console(opts)
        : new winston.transports.File({ filename: '/dev/null' })
    ]

    this.winston = new winston.Logger({
      level,
      levels,
      transports
    })

    return this.winston
  }

  /**
   * @param {string} module     name of the module
   * @param {Object} [options]  winston: a winston logger instance, context: { ... }
   */
  constructor (module, options) {
    options = options || {}
    this.winston = options.winston || Logger.getWinstonLogger()
    this.context = Object.assign({ pid: process.pid }, Logger.context, options.context, module ? { module } : {})
  }

  _log (level, args) {
    let values = []
    let context = Object.assign({}, this.context)

    // Make sure all objects are logged as context
    for (let arg of args) {
      if (arg instanceof Error) {
        values.push(arg.message)
        Object.assign(context, { error: arg })
      } else if (_.isObject(arg)) {
        Object.assign(context, arg)
      } else {
        values.push(arg)
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

module.exports = Logger
