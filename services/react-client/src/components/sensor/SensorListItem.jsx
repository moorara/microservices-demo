import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'

// Presentational Component
const SensorListItem = ({ sensor }) => (
  <tr>
    <td><Link to={`/sensors/${sensor.id}`}>{sensor.name}</Link></td>
    <td>{sensor.unit}</td>
    <td>{sensor.minSafe}</td>
    <td>{sensor.maxSafe}</td>
  </tr>
)

SensorListItem.propTypes = {
  sensor: PropTypes.object.isRequired
}

export default SensorListItem
