/* eslint-env mocha */
const sinon = require('sinon')

const Middleware = require('../../../middleware')

describe('Middleware', () => {
  let req, _req
  let res, _res

  beforeEach(() => {
    req = {
      is () {}
    }
    _req = sinon.mock(req)

    res = {
      end () {},
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

  describe('normalize', () => {
    let middleware

    beforeEach(() => {
      middleware = Middleware.normalize()
    })

    it('should normalize a request without params or query', done => {
      req.url = '/path'
      middleware(req, res, () => {
        res.statusCode = 200
        res.end()
        done()
      })
    })
    it('should normalize a request with params but no query', done => {
      req.url = '/path/1234'
      req.params = { id: '1234' }
      middleware(req, res, () => {
        res.statusCode = 200
        res.end()
        done()
      })
    })
    it('should normalize a request with params and query', done => {
      req.originalUrl = '/path/1234?name=me'
      req.url = '/1234?query=search'
      req.params = { id: '1234' }
      req.query = { name: 'me' }
      middleware(req, res, () => {
        res.statusCode = 200
        res.end()
        done()
      })
    })
  })

  describe('ensureJson', () => {
    let middleware

    beforeEach(() => {
      middleware = Middleware.ensureJson()
    })

    it('should respond with error when body is invalid json', done => {
      _req.expects('is').withArgs('json').returns(false)
      _res.expects('sendStatus').withArgs(415).returns()
      middleware(req, res)
      _req.verify()
      _res.verify()
      done()
    })
    it('should forward the request along the chain when body is valid json', done => {
      _req.expects('is').withArgs('json').returns(true)
      middleware(req, res, () => {
        _req.verify()
        done()
      })
    })
  })

  describe('catchError', () => {
    let err, middleware

    it('should ignore if response is already sent', done => {
      res.headersSent = true
      middleware = Middleware.catchError()
      middleware(err, req, res)
      _res.verify()
      done()
    })
    it('should error with status 500 and no body when environment is production', done => {
      middleware = Middleware.catchError({ environment: 'production' })
      _res.expects('sendStatus').withArgs(500).returns()
      middleware(err, req, res)
      _res.verify()
      done()
    })
    it('should error with status 500 and error in body when environment is not production', done => {
      err = { message: 'error' }
      middleware = Middleware.catchError({ environment: 'development' })
      _res.expects('status').withArgs(500).returns(res)
      _res.expects('send').withArgs(err).returns()
      middleware(err, req, res)
      _res.verify()
      done()
    })
  })
})
