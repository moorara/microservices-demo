/* eslint-env mocha */
const should = require('should')

const Logger = require('../../../util/logger')

describe('Logger', () => {
  describe('addContext', () => {
    afterEach(() => {
      Logger.context = {}
    })

    it('should add properties to context', () => {
      Logger.addContext({ account: 'dev', environment: 'test' })
      Logger.context.account.should.equal('dev')
      Logger.context.environment.should.equal('test')
    })
  })

  describe('getWinstonLogger', () => {
    const origNodeEnv = process.env.NODE_ENV
    const origLogLevel = process.env.LOG_LEVEL

    afterEach(() => {
      delete Logger.winston
      delete process.env.NODE_ENV
      delete process.env.LOG_LEVEL
    })

    after(() => {
      if (origNodeEnv) process.env.NODE_ENV = origNodeEnv
      if (origLogLevel) process.env.LOG_LEVEL = origLogLevel
    })

    it('should return the existing winston logger if already created', () => {
      Logger.winston = { log () {} }
      const winstonLogger = Logger.getWinstonLogger()
      winstonLogger.should.eql(Logger.winston)
    })
    it('should create a new winston logger with defaults', () => {
      const winstonLogger = Logger.getWinstonLogger()
      should.exist(winstonLogger)
      winstonLogger.level.should.equal('info')
    })
    it('should create a new winston logger with debug log level', () => {
      process.env.LOG_LEVEL = 'debug'
      const winstonLogger = Logger.getWinstonLogger()
      should.exist(winstonLogger)
      winstonLogger.level.should.equal('debug')
    })
    it('should create a new winston logger for production', () => {
      process.env.NODE_ENV = 'production'
      const winstonLogger = Logger.getWinstonLogger()
      should.exist(winstonLogger)
      winstonLogger.level.should.equal('info')
    })
    it('should create a new winston logger for development', () => {
      process.env.NODE_ENV = 'development'
      const winstonLogger = Logger.getWinstonLogger()
      should.exist(winstonLogger)
      winstonLogger.level.should.equal('info')
    })
    it('should create a new winston logger for test', () => {
      process.env.NODE_ENV = 'test'
      const winstonLogger = Logger.getWinstonLogger()
      should.exist(winstonLogger)
      winstonLogger.level.should.equal('info')
    })
  })

  describe('constructor', () => {
    afterEach(() => {
      Logger.context = {}
      Logger.winston = null
    })

    it('should create a new logger with defaults', () => {
      const logger = new Logger()
      should.exist(logger.context)
      should.exist(logger.winston)
      logger.context.should.have.property('pid')
    })
    it('should create a new logger with provided options', () => {
      const winston = { log () {} }
      const logger = new Logger('server', {
        winston,
        context: {
          environment: 'test'
        }
      })
      should.exist(logger.context)
      should.exist(logger.winston)
      logger.context.should.have.property('pid')
      logger.context.module.should.equal('server')
      logger.context.environment.should.equal('test')
      logger.winston.should.eql(winston)
    })
  })

  describe('logging', () => {
    let logger
    const origLogLevel = process.env.LOG_LEVEL

    afterEach(() => {
      Logger.context = {}
      Logger.winston = null
      delete process.env.LOG_LEVEL
    })

    after(() => {
      if (origLogLevel) process.env.LOG_LEVEL = origLogLevel
    })

    describe('trace', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'trace'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in trace level', () => {
        logger.trace('hello, world!')
      })
      it('should log strings and context in trace level', () => {
        logger.trace('hello,', 'world')
      })
      it('should log strings with interpolation and context in trace level', () => {
        logger.trace('hello, %s!', 'world')
      })
      it('should log string, error, and context in trace level', () => {
        const err = new Error('mock error!')
        logger.trace('error occurred:', err)
      })
      it('should log string and context in trace level', () => {
        const context = { framework: 'mocha' }
        logger.trace('hello, world!', context)
      })
    })

    describe('debug', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'debug'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in debug level', () => {
        logger.debug('hello, world!')
      })
      it('should log strings and context in debug level', () => {
        logger.debug('hello,', 'world')
      })
      it('should log strings with interpolation and context in debug level', () => {
        logger.debug('hello, %s!', 'world')
      })
      it('should log string, error, and context in debug level', () => {
        const err = new Error('mock error!')
        logger.debug('error occurred:', err)
      })
      it('should log string and context in debug level', () => {
        const context = { framework: 'mocha' }
        logger.debug('hello, world!', context)
      })
    })

    describe('info', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'info'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in info level', () => {
        logger.info('hello, world!')
      })
      it('should log strings and context in info level', () => {
        logger.info('hello,', 'world')
      })
      it('should log strings with interpolation and context in info level', () => {
        logger.info('hello, %s!', 'world')
      })
      it('should log string, error, and context in info level', () => {
        const err = new Error('mock error!')
        logger.info('error occurred:', err)
      })
      it('should log string and context in info level', () => {
        const context = { framework: 'mocha' }
        logger.info('hello, world!', context)
      })
    })

    describe('warn', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'warn'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in warn level', () => {
        logger.warn('hello, world!')
      })
      it('should log strings and context in warn level', () => {
        logger.warn('hello,', 'world')
      })
      it('should log strings with interpolation and context in warn level', () => {
        logger.warn('hello, %s!', 'world')
      })
      it('should log string, error, and context in warn level', () => {
        const err = new Error('mock error!')
        logger.warn('error occurred:', err)
      })
      it('should log string and context in warn level', () => {
        const context = { framework: 'mocha' }
        logger.warn('hello, world!', context)
      })
    })

    describe('error', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'error'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in error level', () => {
        logger.error('hello, world!')
      })
      it('should log strings and context in error level', () => {
        logger.error('hello,', 'world')
      })
      it('should log strings with interpolation and context in error level', () => {
        logger.error('hello, %s!', 'world')
      })
      it('should log string, error, and context in error level', () => {
        const err = new Error('mock error!')
        logger.error('error occurred:', err)
      })
      it('should log string and context in error level', () => {
        const context = { framework: 'mocha' }
        logger.error('hello, world!', context)
      })
    })

    describe('fatal', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'fatal'
        Logger.addContext({ environment: 'test' })
        logger = new Logger()
      })

      it('should log string and context in fatal level', () => {
        logger.fatal('hello, world!')
      })
      it('should log strings and context in fatal level', () => {
        logger.fatal('hello,', 'world')
      })
      it('should log strings with interpolation and context in fatal level', () => {
        logger.fatal('hello, %s!', 'world')
      })
      it('should log string, error, and context in fatal level', () => {
        const err = new Error('mock error!')
        logger.fatal('error occurred:', err)
      })
      it('should log string and context in fatal level', () => {
        const context = { framework: 'mocha' }
        logger.fatal('hello, world!', context)
      })
    })
  })
})
