import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import SensorListItem from '../SensorListItem'

describe('SensorListItem', () => {
  const props = {
    sensor: {}
  }

  test('renders react component with props', () => {
    props.sensor = { id: 'aaaa-aaaa', siteId: '1111-1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
    const wrapper = shallow(<SensorListItem {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find(`Link [to="/sensors/${props.sensor.id}"]`))
    expect(wrapper.find('td').at(1).text()).toBe(`${props.sensor.unit}`)
    expect(wrapper.find('td').at(2).text()).toBe(`${props.sensor.minSafe}`)
    expect(wrapper.find('td').at(3).text()).toBe(`${props.sensor.maxSafe}`)
  })
  test('renders react component with props', () => {
    props.sesnor = { id: 'bbbb-bbbb', siteId: '2222-2222', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
    const wrapper = shallow(<SensorListItem {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find(`Link [to="/sensors/${props.sensor.id}"]`))
    expect(wrapper.find('td').at(1).text()).toBe(`${props.sensor.unit}`)
    expect(wrapper.find('td').at(2).text()).toBe(`${props.sensor.minSafe}`)
    expect(wrapper.find('td').at(3).text()).toBe(`${props.sensor.maxSafe}`)
  })
})
