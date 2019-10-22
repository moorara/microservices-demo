import * as actions from '../actions/sensor'

export const defaultState = {
  callsInProgress: 0,
  items: []
}

export default function (state = defaultState, action) {
  switch (action.type) {
    case actions.SITE_SENSORS_REQUESTED:
      return {
        ...state,
        callsInProgress: state.callsInProgress + 1
      }

    case actions.SITE_SENSORS_FAILED:
      return {
        ...state,
        callsInProgress: state.callsInProgress - 1
      }

    case actions.SITE_SENSORS_RECEIVED:
      return {
        items: action.sensors,
        callsInProgress: state.callsInProgress - 1
      }

    default:
      return state
  }
}
