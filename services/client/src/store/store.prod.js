import thunk from 'redux-thunk'
import { createStore, applyMiddleware } from 'redux'

import rootReducer, { defaultState } from '../reducers'

/**
 * Creates a Redux store.
 * @param {objecy} options
 *   @param {array} sites preloaded array of site objects
 *   @param {array} sensors preloaded array of sensor objects
 * @return {object} Redux store
 */
export default (options) => {
  options = options || {}

  const preloadedState = { ...defaultState }
  preloadedState.site.items = options.sites || preloadedState.site.items
  preloadedState.sensor.items = options.sensors || preloadedState.sensor.items

  return createStore(
    rootReducer,
    preloadedState,
    applyMiddleware(thunk)
  )
}
