import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import { BrowserRouter as Router } from 'react-router-dom'
import { hot } from 'react-hot-loader'

import Header from './common/Header'
import Routes from './Routes'

// Container Component
const App = (props) => (
  <Router>
    <div>
      <Header isLoading={props.sitesLoading || props.sensorsLoading} />
      <Routes />
    </div>
  </Router>
)

App.propTypes = {
  sitesLoading: PropTypes.bool.isRequired,
  sensorsLoading: PropTypes.bool.isRequired
}

/*
 * Map store state to component props
 * See https://github.com/reactjs/react-redux/blob/master/docs/api.md
 */
const mapStateToProps = (state) => ({
  sitesLoading: state.site.callsInProgress > 0,
  sensorsLoading: state.sensor.callsInProgress > 0,
})

export default hot(module)(connect(mapStateToProps)(App))
