import { combineReducers } from 'redux'

import siteReducer, { siteInitialState } from './site'
import sensorReducer, { sensorInitialState } from './sensor'

export default combineReducers({
  site: siteReducer,
  sensor: sensorReducer,
})

export const initialState = {
  site: siteInitialState,
  sensor: sensorInitialState
}
