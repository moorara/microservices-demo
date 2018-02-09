import React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from './home/HomePage'
import SitePage from './site/SitePage'
import AboutPage from './about/AboutPage'

const Routes = (props) => (
  <Switch>
    <Route exact path="/" component={HomePage} />
    <Route exact path="/sites" component={SitePage} />
    <Route exact path="/about" component={AboutPage} />
    <Redirect from="/about-us" to="/about" />
  </Switch>
)

export default Routes
