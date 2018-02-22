import './styles/styles.css'
import '../node_modules/bulma/css/bulma.css'

import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'

import App from './components/App'
import createStore from './store'

const store = createStore({
  logger: true
})

ReactDOM.render((
  <Provider store={store}>
    <App />
  </Provider>
), document.getElementById('app'))
