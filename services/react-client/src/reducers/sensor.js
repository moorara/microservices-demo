import * as actions from '../actions/sensor'

export const sensorInitialState = {
  callsInProgress: 0,
  items: []
}

export default function (state = sensorInitialState, action) {
  switch (action.type) {
    case actions.SITE_SENSORS_REQUESTED:
      return {
        callsInProgress: state.callsInProgress + 1,
        items: Object.assign([], state.items)
      }

    case actions.SITE_SENSORS_FAILED:
      return {
        callsInProgress: state.callsInProgress - 1,
        items: Object.assign([], state.items)
      }

    case actions.SITE_SENSORS_RECEIVED:
      return {
        callsInProgress: state.callsInProgress - 1,
        items: Object.assign([], action.sensors)
      }

    default:
      return state
  }
}
