import * as actions from '../actions/site'

export const siteInitialState = {
  callsInProgress: 0,
  items: []
}

export default function (state = siteInitialState, action) {
  switch (action.type) {
    case actions.ALL_SITES_REQUESTED:
      return {
        callsInProgress: state.callsInProgress + 1,
        items: Object.assign([], state.items)
      }

    case actions.ALL_SITES_FAILED:
      return {
        callsInProgress: state.callsInProgress - 1,
        items: Object.assign([], state.items)
      }

    case actions.ALL_SITES_RECEIVED:
      return {
        callsInProgress: state.callsInProgress - 1,
        items: Object.assign([], action.sites)
      }

    default:
      return state
  }
}
