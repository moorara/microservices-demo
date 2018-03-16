import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import TextInput from '../TextInput'

describe('TextInput', () => {
  let handleChange

  beforeEach(() => {
    // jest.fn(...) is a shorthand for jest.fn().mockImplementation(...)
    handleChange = jest.fn(e => {})
  })

  test('renders react component with placeholder', () => {
    const wrapper = shallow(<TextInput name="firstName" label="First Name" placeholder="Your Name" onChange={handleChange} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.label').text()).toBe('First Name')
    expect(wrapper.find('div .control').exists()).toBeTruthy()
    expect(wrapper.find('input [type="text"] [name="firstName"] [placeholder="Your Name"] [value=""]').exists()).toBeTruthy()
    expect(wrapper.find('span .icon').exists()).toBeFalsy()
    expect(wrapper.find('p .help').exists()).toBeFalsy()

    wrapper.find('input [type="text"]').simulate('change', 'my name')
    expect(handleChange).toHaveBeenCalledWith('my name')
  })
  test('renders react component with value', () => {
    const wrapper = shallow(<TextInput name="firstName" label="First Name" value="Me" onChange={handleChange} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.label').text()).toBe('First Name')
    expect(wrapper.find('div .control').exists()).toBeTruthy()
    expect(wrapper.find('input [type="text"] [name="firstName"] [placeholder=""] [value="Me"]').exists()).toBeTruthy()
    expect(wrapper.find('span .icon').exists()).toBeFalsy()
    expect(wrapper.find('p .help').exists()).toBeFalsy()

    wrapper.find('input [type="text"]').simulate('change', 'my name')
    expect(handleChange).toHaveBeenCalledWith('my name')
  })
  test('renders react component in error case', () => {
    const wrapper = shallow(<TextInput name="firstName" label="First Name" isError help="name is invalid" onChange={handleChange} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.label').text()).toBe('First Name')
    expect(wrapper.find('div .control').exists()).toBeTruthy()
    expect(wrapper.find('input [type="text"] [name="firstName"] [placeholder=""] [value=""]').exists()).toBeTruthy()
    expect(wrapper.find('span .icon i').prop('className')).toContain('fa-exclamation-circle')
    expect(wrapper.find('p .help').text()).toBe('name is invalid')

    wrapper.find('input [type="text"]').simulate('change', 'my name')
    expect(handleChange).toHaveBeenCalledWith('my name')
  })
  test('renders react component in success case', () => {
    const wrapper = shallow(<TextInput name="firstName" label="First Name" isOk help="name is valid" onChange={handleChange} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('.label').text()).toBe('First Name')
    expect(wrapper.find('div .control').exists()).toBeTruthy()
    expect(wrapper.find('input [type="text"] [name="firstName"] [placeholder=""] [value=""]').exists()).toBeTruthy()
    expect(wrapper.find('span .icon i').prop('className')).toContain('fa-check-circle')
    expect(wrapper.find('p .help').text()).toBe('name is valid')

    wrapper.find('input [type="text"]').simulate('change', 'my name')
    expect(handleChange).toHaveBeenCalledWith('my name')
  })
})
