import './styles/styles.css'
import 'bulma/css/bulma.css'
import 'font-awesome/css/font-awesome.css'

import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'

import App from './components/App'
import createStore from './store'
import { getAllSites } from './actions/site'

const store = createStore({
  logger: true
})

store.dispatch(getAllSites())

ReactDOM.render((
  <Provider store={store}>
    <App />
  </Provider>
), document.getElementById('app'))
