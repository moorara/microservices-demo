/*
 * This is a middleware for server-side rendering of client application.
 */

import 'colors'
import rp from 'request-promise'
import React from 'react'
import { Provider } from 'react-redux'
import { renderToString } from 'react-dom/server'

import createStore from '../../src/store'
import App from '../../src/components/App'

const apiUrl = process.env.API_URL || 'http://localhost:4000'

export default async function () {
  let sites, sensors

  // Load state resources
  try {
    sites = await rp.get(`${apiUrl}/api/v1/sites`)
    sensors = await rp.get(`${apiUrl}/api/v1/sensors`)
  } catch (err) {
    console.log(err.red)
  }

  // Create a Redux store
  const store = createStore({ sites, sensors })

  // Render the component to a string
  const serverSideRenderedHtml = renderToString(
    <Provider store={store}>
      <App />
    </Provider>
  )

  // Send the rendered page back to the client
  return {
    serverSideRenderedHtml,
    preloaded: { sites, sensors }
  }
}
