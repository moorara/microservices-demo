import React from 'react'
import PropTypes from 'prop-types'

import SiteListItem from './SiteListItem'

// Presentational Component
const SiteList = ({ sites }) => (
  <table className="table is-bordered is-striped is-fullwidth is-hoverable">
    <thead>
      <tr>
        <th>Name</th>
        <th>Location</th>
        <th>Priority</th>
        <th>Tags</th>
      </tr>
    </thead>
    <tbody>
      {sites.map(site => <SiteListItem key={site.id} site={site} />)}
    </tbody>
  </table>
)

SiteList.propTypes = {
  sites: PropTypes.array.isRequired
}

export default SiteList
