import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import Loading from '../Loading'

describe('Loading', () => {
  const props = {
    dots: 10,
    interval: 100
  }

  beforeEach(() => {
    // Calling this also resets previous mocks
    jest.useFakeTimers()
  })

  test('renders react component with default props', () => {
    const wrapper = shallow(<Loading />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.dots').text()).toBe('•')

    // Test componentDidMount lifecycle and setInterval methods
    expect(setInterval).toHaveBeenCalledTimes(1)
    expect(setInterval).toHaveBeenLastCalledWith(expect.any(Function), 200)

    // Test callback function passed to setInterval
    expect(wrapper.state('frame')).toBe(1)
    const callback = setInterval.mock.calls[0][0]
    callback()
    expect(wrapper.state('frame')).toBe(2)

    // Test componentWillUnmount lifecycle and clearInterval methods
    wrapper.unmount()
    expect(clearInterval).toHaveBeenCalledTimes(1)
  })
  test('renders react component with custom props', () => {
    const wrapper = shallow(<Loading {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.dots').text()).toBe('•')

    // Test componentDidMount lifecycle and setInterval methods
    expect(setInterval).toHaveBeenCalledTimes(1)
    expect(setInterval).toHaveBeenLastCalledWith(expect.any(Function), 100)

    // Test callback function passed to setInterval
    expect(wrapper.state('frame')).toBe(1)
    const callback = setInterval.mock.calls[0][0]
    callback()
    expect(wrapper.state('frame')).toBe(2)

    // Test componentWillUnmount lifecycle and clearInterval methods
    wrapper.unmount()
    expect(clearInterval).toHaveBeenCalledTimes(1)
  })
})
