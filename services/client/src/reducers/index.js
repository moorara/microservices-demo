import { combineReducers } from 'redux'

import siteReducer, { defaultState as siteDefaultState } from './site'
import sensorReducer, { defaultState as sensorDefaultState } from './sensor'

export const defaultState = {
  site: siteDefaultState,
  sensor: sensorDefaultState
}

export default combineReducers({
  site: siteReducer,
  sensor: sensorReducer,
})
