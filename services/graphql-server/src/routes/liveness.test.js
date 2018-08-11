/* eslint-env mocha */
const express = require('express')
const supertest = require('supertest')

const LivenessRouter = require('./liveness')

describe('LivenessRouter', () => {
  let app, request

  beforeEach(() => {
    app = express()
    app.use(LivenessRouter)
    request = supertest(app)
  })

  it('should return 200', done => {
    request.get('/liveness')
      .expect(200, {}, done)
  })
})
