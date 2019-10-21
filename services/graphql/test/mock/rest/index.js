const chalk = require('chalk')
const express = require('express')
const jsonServer = require('json-server')

const mockData = require('./data')

const app = express()
const port = parseInt(process.env.PORT, 10) || 4400

// Logging
app.use((req, res, next) => {
  const end = res.end
  res.end = (data, encoding) => {
    res.end = end
    res.end(data, encoding)

    const timestamp = chalk.green(new Date().toISOString())
    const status = chalk.black.bgCyan.bold(res.statusCode)
    const method = chalk.black.bgYellow.bold(req.method).padEnd(8)
    const url = chalk.yellow(req.path).padEnd(32)
    console.log(` ${timestamp}  ${status}  ${method} ${url}`)
  }

  next()
})

const apiMiddleware = jsonServer.defaults({ logger: false })
const apiRouter = jsonServer.router(mockData)
app.use('/v1', apiMiddleware, apiRouter)

app.listen(port, err => {
  if (err) {
    return console.log(err)
  }

  console.log(chalk.green(`Mock REST API Server Listening on ${port} ...`))
})
