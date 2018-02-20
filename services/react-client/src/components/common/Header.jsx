import React from 'react'
import { Link } from 'react-router-dom'

const Header = (props) => (
  <section className="hero is-dark is-medium is-bold">
    <div className="hero-body">
      <div className="container">
        <h1 className="title">Control Center Application</h1>
        <h2 className="subtitle">Demo Application using Microservies!</h2>
      </div>
    </div>

    <div className="hero-foot">
      <nav className="tabs">
        <div className="container">
          <ul>
            <li><Link to="/">Home</Link></li>
            <li><Link to="/sites">Sites</Link></li>
            <li><Link to="/about">About</Link></li>
          </ul>
        </div>
      </nav>
    </div>
  </section>
)

export default Header
