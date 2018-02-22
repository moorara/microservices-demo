import { combineReducers } from 'redux'

import siteReducer, { defaultState as siteDefaultState } from './site'
import sensorReducer, { defaultState as sensorDefaultState } from './sensor'

export default combineReducers({
  site: siteReducer,
  sensor: sensorReducer,
})

export const defaultState = {
  site: siteDefaultState,
  sensor: sensorDefaultState
}
