/* eslint-env mocha */
const express = require('express')
const supertest = require('supertest')

const HealthRouter = require('../../../routes/health')

describe('HealthRouter', () => {
  let app, request

  beforeEach(() => {
    app = express()
    app.use('/', HealthRouter)

    request = supertest(app)
  })

  it('should return 200', done => {
    request.get('/')
      .expect(200, {}, done)
  })
})
