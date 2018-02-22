/*
 * This is for serving the client application in production mode.
 */

import path from 'path'
import express from 'express'
import compression from 'compression'

const port = parseInt(process.env.PORT, 10) || 4000

const app = express()

app.use(compression())
app.use(
  express.static(
    path.resolve(__dirname, '../../public')
  )
)

app.listen(port, err => {
  if (err) {
    return console.log(err)
  }
  console.log(`Server listening on http://localhost:${port}`)
})
