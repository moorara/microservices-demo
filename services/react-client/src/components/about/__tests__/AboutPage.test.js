import React from 'react'
import { render } from 'enzyme'
import { renderToJson } from 'enzyme-to-json'

import AboutPage from '../AboutPage'

describe('AboutPage', () => {
  test('renders react component', () => {
    const wrapper = render(<AboutPage />)

    expect(renderToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.title').text()).toBe('About')
    expect(wrapper.find('.subtitle').text()).toBe('Control Center Application')
    expect(wrapper.find('.content').text()).toContain('Go')
    expect(wrapper.find('.content').text()).toContain('React')
    expect(wrapper.find('.content').text()).toContain('Redux')
    expect(wrapper.find('.content').text()).toContain('Docker')
    expect(wrapper.find('.content').text()).toContain('Kubernetes')
    expect(wrapper.find('.content').text()).toContain('Prometheus')
    expect(wrapper.find('.content').text()).toContain('ElasticSearch')
    expect(wrapper.find('.content').text()).toContain('Microservices')
  })
})
