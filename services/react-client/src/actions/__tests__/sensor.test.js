import thunk from 'redux-thunk'
import configureStore from 'redux-mock-store'

const middleware = [ thunk ]
const mockStore = configureStore(middleware)

describe('sensorActions', () => {
  let store

  beforeEach(() => {
    store = mockStore({})
    jest.resetModules()
  })

  describe('getSiteSensors', () => {
    let siteId, sensors

    beforeEach(() => {
      siteId = '1111'
      sensors = [
        { id: 'aaaa', siteId, name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
        { id: 'bbbb', siteId, name: 'temperature', unit: 'farenheit', minSafe: -22, maxSafe: 86 }
      ]
    })

    test('dispatches request action and then failure action', () => {
      const err = new Error('error')
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.reject(err)
      }))

      const sensorActions = require('../sensor')
      const func = sensorActions.getSiteSensors(siteId)
      const dispatch = jest.fn()
      return func(dispatch).catch(e => {
        expect(e).toEqual(err)
        expect(dispatch.mock.calls[0][0]).toEqual({ type: 'SITE_SENSORS_REQUESTED', siteId })
        expect(dispatch.mock.calls[1][0]).toEqual({ type: 'SITE_SENSORS_FAILED', err })
      })
    })
    test('dispatches request action and then success action', () => {
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.resolve(sensors)
      }))

      const sensorActions = require('../sensor')
      const func = sensorActions.getSiteSensors(siteId)
      const dispatch = jest.fn()
      return func(dispatch).catch(() => {
        expect(dispatch.mock.calls[0][0]).toEqual({ type: 'SITE_SENSORS_REQUESTED', siteId })
        expect(dispatch.mock.calls[1][0]).toEqual({ type: 'SITE_SENSORS_RECEIVED', sensors })
      })
    })

    test('dispatches request action and then failure action to mock store', () => {
      const err = new Error('error')
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.reject(err)
      }))

      const sensorActions = require('../sensor')
      const action = sensorActions.getSiteSensors(siteId)
      return store.dispatch(action).then(() => {
        const actions = store.getActions()
        expect(actions[0]).toEqual({ type: 'SITE_SENSORS_REQUESTED', siteId })
        expect(actions[1]).toEqual({ type: 'SITE_SENSORS_FAILED', err })
      })
    })
    test('dispatches request action and then success action to mock store', () => {
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.resolve(sensors)
      }))

      const sensorActions = require('../sensor')
      const action = sensorActions.getSiteSensors(siteId)
      return store.dispatch(action).then(() => {
        const actions = store.getActions()
        expect(actions[0]).toEqual({ type: 'SITE_SENSORS_REQUESTED', siteId })
        expect(actions[1]).toEqual({ type: 'SITE_SENSORS_RECEIVED', sensors })
      })
    })
  })
})
