import thunk from 'redux-thunk'
import { createStore, applyMiddleware } from 'redux'

import rootReducer from '../reducers'

export default (initialState) => createStore(
  rootReducer,
  initialState,
  applyMiddleware(thunk)
)
