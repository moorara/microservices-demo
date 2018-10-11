/* eslint-env mocha */
const should = require('should')

const Logger = require('./logger')

describe('Logger', () => {
  describe('addMetadata', () => {
    it('should add properties to metadata', () => {
      Logger.addMetadata({ account: 'dev', environment: 'test' })
      Logger.metadata.account.should.equal('dev')
      Logger.metadata.environment.should.equal('test')
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
      Logger.winston = null
    })

    it('should create a new logger with defaults', () => {
      const logger = new Logger()
      should.exist(logger.metadata)
      should.exist(logger.winston)
      logger.metadata.should.have.property('pid')
    })
    it('should create a new logger with provided options', () => {
      const winston = { log () {} }
      const logger = new Logger('server', {
        winston,
        metadata: {
          environment: 'test'
        }
      })
      should.exist(logger.metadata)
      should.exist(logger.winston)
      logger.metadata.should.have.property('pid')
      logger.metadata.module.should.equal('server')
      logger.metadata.environment.should.equal('test')
      logger.winston.should.eql(winston)
    })
  })

  describe('logging', () => {
    let logger
    const origLogLevel = process.env.LOG_LEVEL

    afterEach(() => {
      Logger.winston = null
      delete process.env.LOG_LEVEL
    })

    after(() => {
      if (origLogLevel) process.env.LOG_LEVEL = origLogLevel
    })

    describe('trace', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'trace'
        Logger.addMetadata({ environment: 'test' })
        logger = new Logger('logger')
      })

      it('should log string and metadata in trace level', () => {
        logger.trace('hello, world!')
      })
      it('should log strings and metadata in trace level', () => {
        logger.trace('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in trace level', () => {
        const str = 'world'
        logger.trace(`hello, ${str}!`)
      })
      it('should log string, error, and metadata in trace level', () => {
        const err = new Error('mock error!')
        logger.trace('error occurred:', err)
      })
      it('should log string and metadata in trace level', () => {
        const metadata = { framework: 'mocha' }
        logger.trace('hello, world!', metadata)
      })
    })

    describe('debug', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'debug'
        Logger.addMetadata({ environment: 'test' })
        logger = new Logger('logger')
      })

      it('should log string and metadata in debug level', () => {
        logger.debug('hello, world!')
      })
      it('should log strings and metadata in debug level', () => {
        logger.debug('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in debug level', () => {
        logger.debug('hello, %s!', 'world')
      })
      it('should log string, error, and metadata in debug level', () => {
        const err = new Error('mock error!')
        logger.debug('error occurred:', err)
      })
      it('should log string and metadata in debug level', () => {
        const metadata = { framework: 'mocha' }
        logger.debug('hello, world!', metadata)
      })
    })

    describe('info', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'info'
        logger = new Logger('logger')
      })

      it('should log string and metadata in info level', () => {
        logger.info('hello, world!')
      })
      it('should log strings and metadata in info level', () => {
        logger.info('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in info level', () => {
        logger.info('hello, %s!', 'world')
      })
      it('should log string, error, and metadata in info level', () => {
        const err = new Error('mock error!')
        logger.info('error occurred:', err)
      })
      it('should log string and metadata in info level', () => {
        const metadata = { framework: 'mocha' }
        logger.info('hello, world!', metadata)
      })
    })

    describe('warn', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'warn'
        Logger.addMetadata({ environment: 'test' })
        logger = new Logger('logger')
      })

      it('should log string and metadata in warn level', () => {
        logger.warn('hello, world!')
      })
      it('should log strings and metadata in warn level', () => {
        logger.warn('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in warn level', () => {
        logger.warn('hello, %s!', 'world')
      })
      it('should log string, error, and metadata in warn level', () => {
        const err = new Error('mock error!')
        logger.warn('error occurred:', err)
      })
      it('should log string and metadata in warn level', () => {
        const metadata = { framework: 'mocha' }
        logger.warn('hello, world!', metadata)
      })
    })

    describe('error', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'error'
        Logger.addMetadata({ environment: 'test' })
        logger = new Logger('logger')
      })

      it('should log string and metadata in error level', () => {
        logger.error('hello, world!')
      })
      it('should log strings and metadata in error level', () => {
        logger.error('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in error level', () => {
        logger.error('hello, %s!', 'world')
      })
      it('should log string, error, and metadata in error level', () => {
        const err = new Error('mock error!')
        logger.error('error occurred:', err)
      })
      it('should log string and metadata in error level', () => {
        const metadata = { framework: 'mocha' }
        logger.error('hello, world!', metadata)
      })
    })

    describe('fatal', () => {
      beforeEach(() => {
        process.env.LOG_LEVEL = 'fatal'
        Logger.addMetadata({ environment: 'test' })
        logger = new Logger('logger')
      })

      it('should log string and metadata in fatal level', () => {
        logger.fatal('hello, world!')
      })
      it('should log strings and metadata in fatal level', () => {
        logger.fatal('hello,', 'world')
      })
      it('should log strings with interpolation and metadata in fatal level', () => {
        logger.fatal('hello, %s!', 'world')
      })
      it('should log string, error, and metadata in fatal level', () => {
        const err = new Error('mock error!')
        logger.fatal('error occurred:', err)
      })
      it('should log string and metadata in fatal level', () => {
        const metadata = { framework: 'mocha' }
        logger.fatal('hello, world!', metadata)
      })
    })
  })
})
