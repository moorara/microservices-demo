function catchError (err, req, res, next) {
  if (!res.headersSent) {
    if (process.env.NODE_ENV === 'production') {
      return res.sendStatus(500)
    }
    res.status(500).send(err)
  }
}

module.exports = { catchError }
