/*
 * This is a proxy server for:
 *   - proxying to client application in development mode,
 *   - and proxying api requests to a mock server
 * 
 * TODO: currently cannot serve /sockjs-node/info
 */

const opn = require('opn')
const express = require('express')
const request = require('request')
const jsonServer = require('json-server')

const mockData = require('./data')

const port = parseInt(process.env.PORT, 10) || 3001
const clientUrl = process.env.CLIENT_URL || 'http://localhost:3000/'

const apiMiddleware = jsonServer.defaults()
const apiRouter = jsonServer.router(mockData)

const app = express()
app.use('/api/v1', apiMiddleware, apiRouter)
app.get('/*', (req, res) => {
  request(`${clientUrl}${req.originalUrl}`).pipe(res)
})

app.listen(port, err => {
  if (err) {
    return console.log(err.red)
  }

  console.log(`Listening on port ${port} ...`)
  opn(`http://localhost:${port}`)
})
