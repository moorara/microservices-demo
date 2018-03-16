import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router-dom'

// Presentational Component
const SiteListItem = ({ site }) => (
  <tr>
    <td><Link to={`/sites/${site.id}`}>{site.name}</Link></td>
    <td>{site.location}</td>
    <td>{site.priority}</td>
    <td>{site.tags.join(' ')}</td>
  </tr>
)

SiteListItem.propTypes = {
  site: PropTypes.object.isRequired
}

export default SiteListItem
