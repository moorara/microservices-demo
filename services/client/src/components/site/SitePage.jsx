import React from 'react'
import PropTypes from 'prop-types'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import SiteList from './SiteList'
import * as siteActions from '../../actions/site'

// Container Component
export class SitePage extends React.Component {
  static propTypes = {
    history: PropTypes.object.isRequired,
    sites: PropTypes.array.isRequired
  }

  constructor (props) {
    super(props)
    this.handleBack = this.handleBack.bind(this)
    this.handleAddSite = this.handleAddSite.bind(this)
  }

  handleBack () {
    this.props.history.goBack()
  }

  handleAddSite () {
    console.log('Add Site clicked!')
  }

  // See https://reactjs.org/docs/react-component.html#shouldcomponentupdate
  shouldComponentUpdate (nextProps, nextState) {
    return true
  }

  render () {
    const { sites } = this.props
    return (
      <section className="hero is-light is-large is-bold">
        <div className="hero-body">
          <div className="container">
            <h1 className="title">All Sites</h1>
            <SiteList sites={sites} />
            <input className="button is-link is-rounded" type="submit" value="Back" onClick={this.handleBack} />
            <span>&nbsp;&nbsp;</span>
            <input className="button is-success is-rounded" type="submit" value="Add Site" onClick={this.handleAddSite} />
          </div>
        </div>
      </section>
    )
  }
}

/*
 * Map store state to component props
 *   ownProps: the props passed to the connected component
 * See https://github.com/reactjs/react-redux/blob/master/docs/api.md
 */
export function mapStateToProps (state, ownProps) {
  return {
    sites: state.site.items
  }
}

/*
 * Wrap each action creator with a dispatch call to store and pass actions to component props
 *   ownProps: the props passed to the connected component
 * See https://github.com/reactjs/react-redux/blob/master/docs/api.md
 * See https://redux.js.org/api-reference/bindactioncreators
 */
export function mapDispatchToProps (dispatch, ownProps) {
  return {
    actions: bindActionCreators(siteActions, dispatch)
    /*
     * actions = {
     *   getAllSites: () => dispatch(siteActions.getAllSites())
     * }
     */
  }
}

const connectToRedux = connect(mapStateToProps, mapDispatchToProps)

export default connectToRedux(withRouter(SitePage))
