import React from 'react'
import { Link } from 'react-router-dom'

const HomePage = (props) => (
  <section className="hero is-light is-large is-bold">
    <div className="hero-body">
      <div className="container">
        <h1 className="title">Welcome!</h1>
        <div className="content">
          <p>You are landed on home page a cool application built using cool technologies both in the front and in the end!</p>
          <p>
            <Link to="/sites" className="button is-primary is-medium is-outlined">Sites</Link>
            <span>&nbsp;&nbsp;</span>
            <Link to="/about-us" className="button is-info is-medium is-outlined">About</Link>
          </p>
        </div>
      </div>
    </div>
  </section>
)

export default HomePage
