/*
 * This is for serving the client application in production mode.
 */

const path = require('path')
const express = require('express')
const compression = require('compression')

const port = parseInt(process.env.PORT, 10) || 4000

const app = express()

app.use(compression())
app.use(
  express.static(
    path.resolve(__dirname, '../public')
  )
)

// Healthcheck endpoint
app.use('/health', (req, res) => {
  res.sendStatus(200)
})

app.listen(port, err => {
  if (err) {
    return console.log(err)
  }
  console.log(`Server listening on http://localhost:${port}`)
})
