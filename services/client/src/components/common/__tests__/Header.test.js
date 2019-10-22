import React from 'react'
import { Link } from 'react-router-dom'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import Header from '../Header'

describe('Header', () => {
  const props = {
    isLoading: false
  }

  test('renders react component with isLoading false', () => {
    props.isLoading = false
    const wrapper = shallow(<Header {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.title').exists()).toBeTruthy()
    expect(wrapper.find('.subtitle').exists()).toBeTruthy()
    expect(wrapper.find('Loading').exists()).toBeFalsy()
    expect(wrapper.find(Link).find('[to="/"]').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/sites"]').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/about"]').exists()).toBeTruthy()
  })
  test('renders react component with isLoading true', () => {
    props.isLoading = true
    const wrapper = shallow(<Header {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.title').exists()).toBeTruthy()
    expect(wrapper.find('.subtitle').exists()).toBeTruthy()
    expect(wrapper.find('Loading').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/"]').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/sites"]').exists()).toBeTruthy()
    expect(wrapper.find(Link).find('[to="/about"]').exists()).toBeTruthy()
  })
})
