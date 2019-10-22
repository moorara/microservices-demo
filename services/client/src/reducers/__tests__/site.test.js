import siteReducer, { siteInitialState } from '../site'

describe('siteReducer', () => {
  let state

  beforeEach(() => {
    state = Object.assign({}, siteInitialState, {
      callsInProgress: 1
    })
  })

  describe('ALL_SITES', () => {
    test('ALL_SITES_REQUESTED', () => {
      const action = {
        type: 'ALL_SITES_REQUESTED'
      }

      const nextState = siteReducer(state, action)
      expect(nextState.callsInProgress).toBe(2)
      expect(nextState.items).toEqual(state.items)
    })
    test('ALL_SITES_FAILED', () => {
      const action = {
        type: 'ALL_SITES_FAILED',
        err: new Error('error')
      }

      const nextState = siteReducer(state, action)
      expect(nextState.callsInProgress).toBe(0)
      expect(nextState.items).toEqual(state.items)
    })
    test('ALL_SITES_RECEIVED', () => {
      const action = {
        type: 'ALL_SITES_RECEIVED',
        sites: [
          { id: '1111', name: 'Power Plant', location: 'Ottawa, ON', tags: [ 'hydro' ], priority: 3 },
          { id: '2222', name: 'Gas Station', location: 'Ottawa, ON', tags: [ 'fuel' ], priority: 2 },
        ]
      }

      const nextState = siteReducer(state, action)
      expect(nextState.callsInProgress).toBe(0)
      expect(nextState.items).toEqual(action.sites)
    })
  })
})
