// https://redux.js.org/docs/advanced/AsyncActions.html

import sensorApi from '../api/sensor'

export const SITE_SENSORS_REQUESTED = 'SITE_SENSORS_REQUESTED'
export const SITE_SENSORS_RECEIVED = 'SITE_SENSORS_RECEIVED'
export const SITE_SENSORS_FAILED = 'SITE_SENSORS_FAILED'

export function requestSiteSensors (siteId) {
  return {
    type: SITE_SENSORS_REQUESTED,
    siteId
  }
}

export function receivedSiteSensors (sensors) {
  return {
    type: SITE_SENSORS_RECEIVED,
    sensors
  }
}

export function failedSiteSensors (err) {
  return {
    type: SITE_SENSORS_FAILED,
    err
  }
}

export function getSiteSensors (siteId) {
  return dispatch => {
    dispatch(requestSiteSensors(siteId))

    return sensorApi.all(siteId)
      .then(
        sensors => dispatch(receivedSiteSensors(sensors)),
        err => dispatch(failedSiteSensors(err))
      )
  }
}
