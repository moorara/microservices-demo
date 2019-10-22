import sensorReducer, { sensorInitialState } from '../sensor'

describe('sensorReducer', () => {
  let state

  beforeEach(() => {
    state = Object.assign({}, sensorInitialState, {
      callsInProgress: 1
    })
  })

  describe('SITE_SENSORS', () => {
    test('SITE_SENSORS_REQUESTED', () => {
      const action = {
        type: 'SITE_SENSORS_REQUESTED',
        siteId: '1111'
      }

      const nextState = sensorReducer(state, action)
      expect(nextState.callsInProgress).toBe(2)
      expect(nextState.items).toEqual(state.items)
    })
    test('SITE_SENSORS_FAILED', () => {
      const action = {
        type: 'SITE_SENSORS_FAILED',
        err: new Error('error')
      }

      const nextState = sensorReducer(state, action)
      expect(nextState.callsInProgress).toBe(0)
      expect(nextState.items).toEqual(state.items)
    })
    test('SITE_SENSORS_RECEIVED', () => {
      const action = {
        type: 'SITE_SENSORS_RECEIVED',
        sensors: [
          { id: 'aaaa', siteId: '1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
          { id: 'bbbb', siteId: '1111', name: 'temperature', unit: 'farenheit', minSafe: -22, maxSafe: 86 }
        ]
      }

      const nextState = sensorReducer(state, action)
      expect(nextState.callsInProgress).toBe(0)
      expect(nextState.items).toEqual(action.sensors)
    })
  })
})
