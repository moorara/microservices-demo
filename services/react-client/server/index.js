/*
 * This is for serving the client application in production mode.
 */

const path = require('path')
const express = require('express')
const compression = require('compression')
const winston = require('winston')
const expressWinston = require('express-winston')

const port = parseInt(process.env.PORT, 10) || 4000

const app = express()

const logger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  transports: [
    new winston.transports.Console()
  ]
})

const loggerMiddleware = expressWinston.logger({
  winstonInstance: logger,
  statusLevels: true
})

// Healthcheck endpoint
app.get('/health', (req, res) => {
  res.sendStatus(200)
})

// Handling api paths
app.use('/api', (req, res) => {
  res.sendStatus(404)
})

app.use(compression())
app.use(loggerMiddleware)

// Serving assets
app.use(express.static(
  path.join(__dirname, '../build')
))

// Serving application
app.get('/*', (req, res) => {
  res.sendFile(path.join(__dirname, '../build', 'index.html'))
})

app.listen(port, err => {
  if (err) {
    return logger.error(err)
  }

  logger.info(`Listening on port ${port} ...`)
})
