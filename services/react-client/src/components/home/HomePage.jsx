import React from 'react'
import { Link } from 'react-router-dom'

const Home = (props) => (
  <section className="hero is-light is-large is-bold">
    <div className="hero-body">
      <div className="container">
        <h1 className="title">Welcome!</h1>
        <div className="columns">
          <div className="column">
            <Link to="/about-us" className="button is-info is-large is-outlined">About</Link>
          </div>
        </div>
      </div>
    </div>
  </section>
)

export default Home
