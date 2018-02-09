import React from 'react'
import { BrowserRouter as Router } from 'react-router-dom'

import Header from './common/Header'
import Routes from './Routes'

const App = (props) => (
  <Router>
    <div>
      <Header />
      <Routes />
    </div>
  </Router>
)

export default App
