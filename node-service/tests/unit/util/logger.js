/* eslint-env mocha */
const should = require('should')

const Logger = require('../../../util/logger')

describe('Logger', () => {
  describe('constructor', () => {
    let logger

    afterEach(() => {
      Logger.winston = null
    })

    it('should create a new logger with defaults', () => {
      logger = new Logger()
      should.exist(logger.winston)
      should.exist(logger.context)
    })
    it('should create a new logger with custom settings', () => {
      logger = new Logger('test', {
        winston: {},
        context: {
          library: 'mocha'
        }
      })
      should.exist(logger.winston)
      should.exist(logger.context)
      logger.context.library.should.equal('mocha')
    })
    it('should create a new logger for development environment', () => {
      const nodeEnv = process.env.NODE_ENV
      process.env.NODE_ENV = 'development'
      logger = new Logger('test')
      should.exist(logger.winston)
      should.exist(logger.context)
      process.env.NODE_ENV = nodeEnv
    })
  })

  describe('logging', () => {
    let logger

    beforeEach(() => {
      logger = new Logger('test')
    })

    afterEach(() => {
      Logger.winston = null
    })

    it('should log in trace level to std out', () => {
      logger.trace(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log in debug level to std out', () => {
      logger.debug(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log in info level to std out', () => {
      logger.info(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log in warn level to std out', () => {
      logger.warn(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log with an error in warn level to std out', () => {
      logger.warn(
        'Hello, %s!', 'World',
        new Error('error')
      )
    })

    it('should log in error level to std out', () => {
      logger.error(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log with an error in error level to std out', () => {
      logger.error(
        'Hello, %s!', 'World',
        new Error('error')
      )
    })

    it('should log in fatal level to std out', () => {
      logger.fatal(
        'Hello, %s!', 'World', 'Hello', 'Universe',
        { sender: 'moorara' },
        { receivers: [ 'world', 'universe' ] }
      )
    })

    it('should log with an error in fatal level to std out', () => {
      logger.fatal(
        'Hello, %s!', 'World',
        new Error('error')
      )
    })
  })
})
