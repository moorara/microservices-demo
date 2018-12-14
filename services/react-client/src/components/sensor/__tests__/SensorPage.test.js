import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import * as actions from '../../../actions/sensor'
import { SensorPage, mapStateToProps, mapDispatchToProps } from '../SensorPage'

describe('SensorPage', () => {
  describe('SensorPage', () => {
    let nextProps, nextState
    const props = {
      history: {},
      match: { params: {} },
      site: {},
      sensors: [],
      actions: {}
    }

    beforeEach(() => {
      props.history.goBack = jest.fn()
      props.actions.getSiteSensors = jest.fn()
    })

    test('renders connected react component with props', () => {
      props.match.params.id = '1111-1111'
      props.site = { id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] }
      props.sensors = [
        { id: 'aaaa-aaaa', siteId: '1111-1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
      ]

      const wrapper = shallow(<SensorPage {...props} />)

      expect(shallowToJson(wrapper)).toMatchSnapshot()
      expect(wrapper.find('SensorList').prop('sensors')).toEqual(props.sensors)
      expect(wrapper.find('input').find({ value: 'Back' }).exists()).toBeTruthy()
      expect(wrapper.find('input').find({ value: 'Add Sensor' }).exists()).toBeTruthy()

      expect(wrapper.instance().shouldComponentUpdate(nextProps, nextState)).toBe(true)

      wrapper.find('input').find({ value: 'Back' }).simulate('click')
      expect(props.history.goBack).toHaveBeenCalled()
    })
    test('renders connected react component with props', () => {
      props.match.params.id = '2222-2222'
      props.sites = { id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }
      props.sensors = [
        { id: 'bbbb-bbbb', siteId: '2222-2222', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      ]

      const wrapper = shallow(<SensorPage {...props} />)

      expect(shallowToJson(wrapper)).toMatchSnapshot()
      expect(wrapper.find('SensorList').prop('sensors')).toEqual(props.sensors)
      expect(wrapper.find('input').find({ value: 'Back' }).exists()).toBeTruthy()
      expect(wrapper.find('input').find({ value: 'Add Sensor' }).exists()).toBeTruthy()

      expect(wrapper.instance().shouldComponentUpdate(nextProps, nextState)).toBe(true)

      wrapper.find('input').find({ value: 'Back' }).simulate('click')
      expect(props.history.goBack).toHaveBeenCalled()
    })
  })

  describe('mapStateToProps', () => {
    const state = {
      site: {
        items: []
      },
      sensor: {
        items: []
      }
    }
    const ownProps = {
      match: {
        params: {}
      }
    }

    test('maps store state to component props', () => {
      ownProps.match.params.id = '1111-1111'
      state.site.items = [{ id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] }]
      state.sensor.items = [
        { id: 'aaaa-aaaa', siteId: '1111-1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
      ]

      const props = mapStateToProps(state, ownProps)
      expect(props.site).toEqual(state.site.items[0])
      expect(props.sensors).toEqual(state.sensor.items)
    })
    test('maps store state to component props', () => {
      ownProps.match.params.id = '2222-2222'
      state.site.items = [{ id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }]
      state.sensor.items = [
        { id: 'bbbb-bbbb', siteId: '2222-2222', name: 'temperature', unit: 'fahrenheit', minSafe: -22, maxSafe: 86 }
      ]

      const props = mapStateToProps(state, ownProps)
      expect(props.site).toEqual(state.site.items[0])
      expect(props.sensors).toEqual(state.sensor.items)
    })
  })

  describe('mapDispatchToProps', () => {
    let dispatch

    beforeEach(() => {
      dispatch = jest.fn()
      actions.getSiteSensors = jest.fn().mockReturnValue({
        type: 'GET_SITE_SENSORS'
      })
    })

    test('maps store dispatch and actions to component props', () => {
      const props = mapDispatchToProps(dispatch)
      props.actions.getSiteSensors('1111-1111')

      expect(actions.getSiteSensors).toHaveBeenCalledWith('1111-1111')
      expect(dispatch).toHaveBeenCalledWith({
        type: 'GET_SITE_SENSORS'
      })
    })
  })
})
