const url = require('url')

const normalizeReq = (req) => {
  let endpoint = req.originalUrl || req.url
  endpoint = url.parse(endpoint).pathname
  for (let p in req.params) {
    if (req.params.hasOwnProperty(p)) {
      endpoint = endpoint.replace(req.params[p], `:${p}`)
    }
  }
  req.endpoint = endpoint
}

const normalizeRes = (res) => {
  const statusClass = Math.floor(res.statusCode / 100)
  res.statusClass = `${statusClass}xx`
}

module.exports = {
  normalize () {
    return (req, res, next) => {
      var end = res.end
      res.end = function (data, encoding, callback) {
        normalizeReq(req)
        normalizeRes(res)
        res.end = end
        res.end(data, encoding, callback)
      }
      next()
    }
  },

  ensureJson () {
    return (req, res, next) => {
      if (!req.is('json')) {
        return res.sendStatus(415)
      }
      next()
    }
  },

  catchError (options) {
    options = options || {}
    return (err, req, res, next) => {
      if (!res.headersSent) {
        if (options.environment === 'production') {
          return res.sendStatus(500)
        }
        res.status(500).send(err)
      }
    }
  }
}
