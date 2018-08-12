const winston = require('winston')

const levels = { fatal: 0, error: 1, warn: 2, info: 3, debug: 4, trace: 5 }

class Logger {
  static addMetadata (metadata) {
    Object.assign(this.metadata, metadata)
  }

  static getWinstonLogger () {
    if (this.winston) {
      return this.winston
    }

    let format
    const level = process.env.LOG_LEVEL || 'info'

    if (process.env.NODE_ENV === 'development') {
      format = winston.format.combine(
        winston.format.splat(),
        winston.format.timestamp(),
        winston.format.prettyPrint(),
        winston.format.colorize(),
        winston.format.simple()
      )
    } else {
      format = winston.format.combine(
        winston.format.splat(),
        winston.format.timestamp(),
        winston.format.printf(info => JSON.stringify(info.message)),
        winston.format.json()
      )
    }

    const transports = [
      process.env.NODE_ENV === 'test'
        ? new winston.transports.File({ filename: '/dev/null' })
        : new winston.transports.Console({ handleExceptions: true })
    ]

    this.winston = winston.createLogger({
      level,
      levels,
      format,
      transports
    })

    return this.winston
  }

  /**
   * @param {string} module   name of the module
   * @param {Object} options  winston: a winston logger instance, metadata: { ... }
   */
  constructor (module, options) {
    options = options || {}
    this.winston = options.winston || Logger.getWinstonLogger()
    this.metadata = Object.assign({ pid: process.pid }, Logger.metadata, options.metadata, module ? { module } : {})
  }

  _log (level, args) {
    let values = []
    let metadata = Object.assign({}, this.metadata)

    // Make sure all objects are logged as metadata
    for (let arg of args) {
      if (arg instanceof Error) {
        values.push(arg.message)
        Object.assign(metadata, { error: arg })
      } else if (typeof arg === 'object') {
        Object.assign(metadata, arg)
      } else {
        values.push(arg)
      }
    }

    this.winston.log(level, values.join(' '), metadata)
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

Logger.metadata = {}

module.exports = Logger
