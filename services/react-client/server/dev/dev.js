/*
 * This is for serving the client application in development mode.
 */

import 'colors'
import opn from 'opn'
import path from 'path'
import express from 'express'
import webpack from 'webpack'
import webpackDevMiddleware from 'webpack-dev-middleware'
import webpackHotMiddleware from 'webpack-hot-middleware'

import mockApiServer from './api'
import webpackConfig from '../../webpack.dev'

const port = parseInt(process.env.PORT, 10) || 4000

// Update config for hot reloading
const hmr = 'webpack-hot-middleware/client?path=/hmr&timeout=2000&overlay=true'
webpackConfig.entry.app = [ hmr, webpackConfig.entry.app ]

const app = express()
const transpiler = webpack(webpackConfig)

const devMiddleware = webpackDevMiddleware(transpiler, {
  noInfo: true,
  publicPath: webpackConfig.output.publicPath
})

const hotMiddleware = webpackHotMiddleware(transpiler, {
  path: '/hmr',
  heartbeat: 1000
})

// Respond with client on every path
const filename = path.join(transpiler.outputPath, 'index.html')
const allMiddleware = (req, res) => {
  devMiddleware.fileSystem.readFile(filename, (err, file) => {
    err ? res.sendStatus(404) : res.send(file.toString())
  })
}

mockApiServer(app)
app.use(devMiddleware)
app.use(hotMiddleware)
app.use(allMiddleware)

app.listen(port, err => {
  if (err) {
    return console.log(err.red)
  }
  console.log(`Dev server listening on http://localhost:${port}`.bold.green)
  opn(`http://localhost:${port}`)
})
