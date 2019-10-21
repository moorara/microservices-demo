import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'

import Loading from './Loading'

const Header = ({ isLoading }) => (
  <section className="hero is-dark is-medium is-bold">
    <div className="hero-body">
      <div className="container">
        <h1 className="title">Control Center</h1>
        <h2 className="subtitle">Monitoring Remote Sites in Real-Time!</h2>
      </div>
    </div>

    <div className="hero-foot">
      <nav className="tabs">
        <div className="container">
          <ul>
            <li><Link to="/">Home</Link></li>
            <li><Link to="/sites">Sites</Link></li>
            <li><Link to="/about">About</Link></li>
            { isLoading && <Loading dots={5} interval={100} /> }
          </ul>
        </div>
      </nav>
    </div>
  </section>
)

Header.propTypes = {
  isLoading: PropTypes.bool.isRequired
}

export default Header
