import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { createStore, applyMiddleware } from 'redux'

import rootReducer from '../reducers'

export default (initialState, options) => {
  options = options || {}
  const middleware = [ thunk ]
  if (options.logger) {
    middleware.push(createLogger())
  }

  return createStore(
    rootReducer,
    initialState,
    applyMiddleware(...middleware)
  )
}
