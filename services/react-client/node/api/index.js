import jsonServer from 'json-server'

import mockData from './data'

export default function (app) {
  const apiMiddleware = jsonServer.defaults()
  const apiRouter = jsonServer.router(mockData)
  app.use('/api/v1', apiMiddleware, apiRouter)
}
