/* eslint-env mocha */
const sinon = require('sinon')

const Middleware = require('.')

describe('Middleware', () => {
  let req, _req
  let res, _res

  beforeEach(() => {
    req = {}
    _req = sinon.mock(req)

    res = {
      status () {},
      send () {},
      sendStatus () {}
    }
    _res = sinon.mock(res)
  })

  afterEach(() => {
    _req.restore()
    _res.restore()
  })

  describe('catchError', () => {
    let middleware
    const origNodeEnv = process.env.NODE_ENV

    afterEach(() => {
      delete process.env.NODE_ENV
    })

    after(() => {
      process.env.NODE_ENV = origNodeEnv
    })

    it('should ignore if response is already sent', done => {
      res.headersSent = true
      middleware = Middleware.catchError

      middleware(null, req, res)
      done()
    })
    it('should error with status 500 and no body when environment is production', done => {
      process.env.NODE_ENV = 'production'
      middleware = Middleware.catchError
      _res.expects('sendStatus').withArgs(500).returns()

      middleware(null, req, res)
      _res.verify()
      done()
    })
    it('should error with status 500 and error in body when environment is not production', done => {
      process.env.NODE_ENV = 'development'
      middleware = Middleware.catchError
      const err = new Error('error')
      _res.expects('status').withArgs(500).returns(res)
      _res.expects('send').withArgs(err).returns()

      middleware(err, req, res)
      _res.verify()
      done()
    })
  })
})
