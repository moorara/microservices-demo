import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import SensorList from '../SensorList'

describe('SensorList', () => {
  const props = {
    sensors: []
  }

  test('renders react component with props', () => {
    props.sensors = [
      { id: 'aaaa-aaaa', siteId: '1111-1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
      { id: 'bbbb-bbbb', siteId: '2222-2222', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
    ]
    const wrapper = shallow(<SensorList {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('thead tr th').at(0).text()).toBe('Name')
    expect(wrapper.find('thead tr th').at(1).text()).toBe('Unit')
    expect(wrapper.find('thead tr th').at(2).text()).toBe('Minimum Safe Value')
    expect(wrapper.find('thead tr th').at(3).text()).toBe('Maximum Safe Value')
    expect(wrapper.find('SensorListItem').get(0).props.sensor).toEqual(props.sensors[0])
    expect(wrapper.find('SensorListItem').get(1).props.sensor).toEqual(props.sensors[1])
  })
  test('renders react component with props', () => {
    props.sensors = [
      { id: 'cccc-cccc', siteId: '3333-3333', name: 'pressure', unit: 'pascal', minSafe: 50000, maxSafe: 100000 },
      { id: 'dddd-dddd', siteId: '4444-4444', name: 'pressure', unit: 'atmosphere', minSafe: 0.5, maxSafe: 1 }
    ]
    const wrapper = shallow(<SensorList {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('thead tr th').at(0).text()).toBe('Name')
    expect(wrapper.find('thead tr th').at(1).text()).toBe('Unit')
    expect(wrapper.find('thead tr th').at(2).text()).toBe('Minimum Safe Value')
    expect(wrapper.find('thead tr th').at(3).text()).toBe('Maximum Safe Value')
    expect(wrapper.find('SensorListItem').get(0).props.sensor).toEqual(props.sensors[0])
    expect(wrapper.find('SensorListItem').get(1).props.sensor).toEqual(props.sensors[1])
  })
})
