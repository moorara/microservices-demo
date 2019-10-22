import React from 'react'
import { Link } from 'react-router-dom'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import HomePage from '../HomePage'

describe('HomePage', () => {
  test('renders react component', () => {
    const wrapper = shallow(<HomePage />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.title').exists()).toBeTruthy()
    expect(wrapper.find('.content').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/sites"]').children().text()).toBe('Sites')
    expect(wrapper.find(Link).find('[to="/about-us"]').children().text()).toBe('About')
  })
})
