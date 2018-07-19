import React from 'react'
import PropTypes from 'prop-types'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { withRouter } from 'react-router-dom'

import SensorList from './SensorList'
import * as sensorActions from '../../actions/sensor'

// Container Component
export class SensorPage extends React.Component {
  static propTypes = {
    history: PropTypes.object.isRequired,
    match: PropTypes.object.isRequired,
    site: PropTypes.object.isRequired,
    sensors: PropTypes.array.isRequired,
    actions: PropTypes.object.isRequired,
  }

  constructor (props) {
    super(props)
    this.handleBack = this.handleBack.bind(this)
    this.handleAddSensor = this.handleAddSensor.bind(this)
  }

  handleBack () {
    this.props.history.goBack()
  }

  handleAddSensor () {
    console.log('Add Sensor clicked!')
  }

  // See https://reactjs.org/docs/react-component.html#shouldcomponentupdate
  shouldComponentUpdate (nextProps, nextState) {
    return true
  }

  componentDidMount () {
    const { actions } = this.props
    const siteId = this.props.match.params.id
    actions.getSiteSensors(siteId)
  }

  render () {
    const { site, sensors } = this.props
    return (
      <section className="hero is-light is-large is-bold">
        <div className="hero-body">
          <div className="container">
            <h1 className="title">{`${site.name} Sensors`}</h1>
            <SensorList sensors={sensors} />
            <input className="button is-link is-rounded" type="submit" value="Back" onClick={this.handleBack} />
            <span>&nbsp;&nbsp;</span>
            <input className="button is-danger is-rounded" type="submit" value="Add Sensor" onClick={this.handleAddSensor} />
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
  const siteId = ownProps.match.params.id
  return {
    site: state.site.items.find(site => site.id === siteId) || {},
    sensors: state.sensor.items.filter(sensor => sensor.siteId === siteId)
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
    actions: bindActionCreators(sensorActions, dispatch)
    /*
     * actions = {
     *   getSiteSensors: siteId => dispatch(sensorActions.getSiteSensors(siteId))
     * }
     */
  }
}

const connectToRedux = connect(mapStateToProps, mapDispatchToProps)

export default connectToRedux(withRouter(SensorPage))
