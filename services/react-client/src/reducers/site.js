import * as actions from '../actions/site'

export const defaultState = {
  callsInProgress: 0,
  items: []
}

export default function (state = defaultState, action) {
  switch (action.type) {
    case actions.ALL_SITES_REQUESTED:
      return {
        ...state,
        callsInProgress: state.callsInProgress + 1
      }

    case actions.ALL_SITES_FAILED:
      return {
        ...state,
        callsInProgress: state.callsInProgress - 1
      }

    case actions.ALL_SITES_RECEIVED:
      return {
        items: action.sites,
        callsInProgress: state.callsInProgress - 1,
      }

    default:
      return state
  }
}
