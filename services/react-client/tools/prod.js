import path from 'path'
import express from 'express'
import compression from 'compression'

const port = process.env.PORT || 4000
const app = express()

app.use(compression())
app.use(
  express.static(
    path.resolve(__dirname, '../public')
  )
)

app.listen(port, err => {
  if (err) {
    console.log(err)
    return 1
  }
  console.log(`Server listening on http://localhost:${port}`)
})
