// https://redux.js.org/docs/advanced/AsyncActions.html

import siteApi from '../api/site'

export const ALL_SITES_REQUESTED = 'ALL_SITES_REQUESTED'
export const ALL_SITES_RECEIVED = 'ALL_SITES_RECEIVED'
export const ALL_SITES_FAILED = 'ALL_SITES_FAILED'

export function requestAllSites () {
  return {
    type: ALL_SITES_REQUESTED
  }
}

export function receivedAllSites (sites) {
  return {
    type: ALL_SITES_RECEIVED,
    sites
  }
}

export function failedAllSites (err) {
  return {
    type: ALL_SITES_FAILED,
    err
  }
}

export function getAllSites () {
  return dispatch => {
    dispatch(requestAllSites())

    return siteApi.all()
      .then(
        sites => dispatch(receivedAllSites(sites)),
        err => dispatch(failedAllSites(err))
      )
  }
}
