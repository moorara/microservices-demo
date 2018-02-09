import createStore from '../../store'
import { initialState } from '../../reducers'

describe('store', () => {
  let store

  beforeEach(() => {
    store = createStore(initialState)
    jest.resetModules()
  })

  describe('getAllSites', () => {
    let sites

    beforeEach(() => {
      sites = [
        { id: '1111', name: 'Power Plant', location: 'Ottawa, ON', tags: [ 'hydro' ], priority: 3 },
        { id: '2222', name: 'Gas Station', location: 'Ottawa, ON', tags: [ 'fuel' ], priority: 2 },
      ]
    })

    test('dispatches request action and then failure action to store', () => {
      const err = new Error('error')
      jest.doMock('../../api/site', () => ({
        all: () => Promise.reject(err)
      }))

      const siteActions = require('../../actions/site')
      const action = siteActions.getAllSites()
      return store.dispatch(action).then(() => {
        const state = store.getState()
        expect(state).toEqual({
          site: { callsInProgress: 0, items: [] },
          sensor: { callsInProgress: 0, items: [] }
        })
      })
    })
    test('dispatches request action and then success action to store', () => {
      jest.doMock('../../api/site', () => ({
        all: () => Promise.resolve(sites)
      }))

      const siteActions = require('../../actions/site')
      const action = siteActions.getAllSites()
      return store.dispatch(action).then(() => {
        const state = store.getState()
        expect(state).toEqual({
          site: { callsInProgress: 0, items: sites },
          sensor: { callsInProgress: 0, items: [] }
        })
      })
    })
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

    test('dispatches request action and then failure action to store', () => {
      const err = new Error('error')
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.reject(err)
      }))

      const sensorActions = require('../../actions/sensor')
      const action = sensorActions.getSiteSensors(siteId)
      return store.dispatch(action).then(() => {
        const state = store.getState()
        expect(state).toEqual({
          site: { callsInProgress: 0, items: [] },
          sensor: { callsInProgress: 0, items: [] }
        })
      })
    })
    test('dispatches request action and then success action to store', () => {
      jest.doMock('../../api/sensor', () => ({
        all: () => Promise.resolve(sensors)
      }))

      const sensorActions = require('../../actions/sensor')
      const action = sensorActions.getSiteSensors(siteId)
      return store.dispatch(action).then(() => {
        const state = store.getState()
        expect(state).toEqual({
          site: { callsInProgress: 0, items: [] },
          sensor: { callsInProgress: 0, items: sensors }
        })
      })
    })
  })
})
