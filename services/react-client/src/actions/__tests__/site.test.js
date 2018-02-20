import thunk from 'redux-thunk'
import configureStore from 'redux-mock-store'

const middleware = [ thunk ]
const mockStore = configureStore(middleware)

describe('siteActions', () => {
  let store

  beforeEach(() => {
    store = mockStore({})
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

    test('dispatches request action and then failure action', () => {
      const err = new Error('error')
      jest.doMock('../../api/site', () => ({
        all: () => Promise.reject(err)
      }))

      const siteActions = require('../site')
      const func = siteActions.getAllSites()
      const dispatch = jest.fn()
      return func(dispatch).catch(e => {
        expect(e).toEqual(err)
        expect(dispatch.mock.calls[0][0]).toEqual({ type: 'ALL_SITES_REQUESTED' })
        expect(dispatch.mock.calls[1][0]).toEqual({ type: 'ALL_SITES_FAILED', err })
      })
    })
    test('dispatches request action and then success action', () => {
      jest.doMock('../../api/site', () => ({
        all: () => Promise.resolve(sites)
      }))

      const siteActions = require('../site')
      const func = siteActions.getAllSites()
      const dispatch = jest.fn()
      return func(dispatch).then(() => {
        expect(dispatch.mock.calls[0][0]).toEqual({ type: 'ALL_SITES_REQUESTED' })
        expect(dispatch.mock.calls[1][0]).toEqual({ type: 'ALL_SITES_RECEIVED', sites })
      })
    })

    test('dispatches request action and then failure action to mock store', () => {
      const err = new Error('error')
      jest.doMock('../../api/site', () => ({
        all: () => Promise.reject(err)
      }))

      const siteActions = require('../site')
      const action = siteActions.getAllSites()
      return store.dispatch(action).then(() => {
        const actions = store.getActions()
        expect(actions[0]).toEqual({ type: 'ALL_SITES_REQUESTED' })
        expect(actions[1]).toEqual({ type: 'ALL_SITES_FAILED', err })
      })
    })
    test('dispatches request action and then success action to mock store', () => {
      jest.doMock('../../api/site', () => ({
        all: () => Promise.resolve(sites)
      }))

      const siteActions = require('../site')
      const action = siteActions.getAllSites()
      return store.dispatch(action).then(() => {
        const actions = store.getActions()
        expect(actions[0]).toEqual({ type: 'ALL_SITES_REQUESTED' })
        expect(actions[1]).toEqual({ type: 'ALL_SITES_RECEIVED', sites })
      })
    })
  })
})
