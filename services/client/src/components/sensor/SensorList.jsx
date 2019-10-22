import React from 'react'
import PropTypes from 'prop-types'

import SensorListItem from './SensorListItem'

// Presentational Component
const SensorList = ({ sensors }) => (
  <table className="table is-bordered is-striped is-fullwidth is-hoverable">
    <thead>
      <tr>
        <th>Name</th>
        <th>Unit</th>
        <th>Minimum Safe Value</th>
        <th>Maximum Safe Value</th>
      </tr>
    </thead>
    <tbody>
      {sensors.map(sensor => <SensorListItem key={sensor.id} sensor={sensor} />)}
    </tbody>
  </table>
)

SensorList.propTypes = {
  sensors: PropTypes.array.isRequired
}

export default SensorList
