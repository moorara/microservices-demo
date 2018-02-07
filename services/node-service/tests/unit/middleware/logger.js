/* eslint-env mocha */
const express = require('express')
const supertest = require('supertest')

const LoggerMiddleware = require('../../../middleware/logger')

describe('LoggerMiddleware', () => {
  let middleware
  let app, request

  describe('http', () => {
    it('should create a logger middleware with defaults', done => {
      middleware = LoggerMiddleware.http()

      app = express()
      app.use(middleware)
      app.use('/:id', (req, res) => res.sendStatus(200))

      request = supertest(app)
      request.get('/1234')
        .expect(200, done)
    })
    it('should create a logger middleware with options', done => {
      middleware = LoggerMiddleware.http({
        level: 'debug',
        skip: (req, res) => false,
        ignoreRoute: (req, res) => false,
        ignoredRoutes: [ '/health', '/metrics' ],
        requestFilter: (req, propName) => req[propName],
        requestWhitelist: [ 'method', 'url' ],
        responseFilter: (res, propName) => res[propName],
        responseWhitelist: [ 'statusCode', 'responseTime' ]
      })

      app = express()
      app.use(middleware)
      app.use('/', (req, res) => res.sendStatus(200))

      request = supertest(app)
      request.get('/')
        .expect(200, done)
    })
  })
})
