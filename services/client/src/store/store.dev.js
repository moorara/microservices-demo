import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { createStore, applyMiddleware } from 'redux'

import rootReducer, { defaultState } from '../reducers'

/**
 * Creates a Redux store.
 * @param {objecy} options
 *   @param {array} sites preloaded array of site objects
 *   @param {array} sensors preloaded array of sensor objects
 *   @param {boolean} logger enables logger middleware
 * @return {object} Redux store
 */
export default (options) => {
  options = options || {}

  const preloadedState = { ...defaultState }
  preloadedState.site.items = options.sites || preloadedState.site.items
  preloadedState.sensor.items = options.sensors || preloadedState.sensor.items

  const middleware = [ thunk ]
  if (options.logger === true) {
    middleware.push(createLogger())
  }

  return createStore(
    rootReducer,
    preloadedState,
    applyMiddleware(...middleware)
  )
}
