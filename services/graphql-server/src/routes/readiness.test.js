/* eslint-env mocha */
const express = require('express')
const supertest = require('supertest')

const ReadinessRouter = require('./readiness')

describe('ReadinessRouter', () => {
  let app, request

  beforeEach(() => {
    app = express()
    app.use(ReadinessRouter)
    request = supertest(app)
  })

  it('should return 200', done => {
    request.get('/readiness')
      .expect(200, {}, done)
  })
})
