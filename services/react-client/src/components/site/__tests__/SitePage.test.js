import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import * as actions from '../../../actions/site'
import { SitePage, mapStateToProps, mapDispatchToProps } from '../SitePage'

describe('SitePage', () => {
  describe('SitePage', () => {
    let nextProps, nextState
    const props = {
      history: {},
      sites: []
    }

    beforeEach(() => {
      props.history.goBack = jest.fn()
    })

    test('renders connected react component with props', () => {
      props.sites = [
        { id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] },
        { id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }
      ]

      const wrapper = shallow(<SitePage {...props} />)

      expect(shallowToJson(wrapper)).toMatchSnapshot()
      expect(wrapper.find('SiteList').prop('sites')).toEqual(props.sites)
      expect(wrapper.find('input [value="Back"]').exists()).toBeTruthy()
      expect(wrapper.find('input [value="Add Site"]').exists()).toBeTruthy()

      expect(wrapper.instance().shouldComponentUpdate(nextProps, nextState)).toBe(true)

      wrapper.find('input [value="Back"]').simulate('click')
      expect(props.history.goBack).toHaveBeenCalled()
    })
    test('renders connected react component with props', () => {
      props.sites = [
        { id: '3333-3333', name: 'Station 3', location: 'Ottawa, ON', priority: 4, tags: [ 'helium', 'nitrogen' ] },
        { id: '4444-4444', name: 'Station 4', location: 'Montreal, QC', priority: 5, tags: [ 'carbon', 'iron' ] }
      ]

      const wrapper = shallow(<SitePage {...props} />)

      expect(shallowToJson(wrapper)).toMatchSnapshot()
      expect(wrapper.find('SiteList').prop('sites')).toEqual(props.sites)
      expect(wrapper.find('input [value="Back"]').exists()).toBeTruthy()
      expect(wrapper.find('input [value="Add Site"]').exists()).toBeTruthy()

      expect(wrapper.instance().shouldComponentUpdate(nextProps, nextState)).toBe(true)

      wrapper.find('input [value="Back"]').simulate('click')
      expect(props.history.goBack).toHaveBeenCalled()
    })
  })

  describe('mapStateToProps', () => {
    const state = {
      site: {
        items: []
      }
    }

    test('maps store state to component props', () => {
      state.site.items = [
        { id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] },
        { id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }
      ]

      const props = mapStateToProps(state)
      expect(props.sites).toEqual(state.site.items)
    })
    test('maps store state to component props', () => {
      state.site.items = [
        { id: '3333-3333', name: 'Station 3', location: 'Ottawa, ON', priority: 4, tags: [ 'helium', 'nitrogen' ] },
        { id: '4444-4444', name: 'Station 4', location: 'Montreal, QC', priority: 5, tags: [ 'carbon', 'iron' ] }
      ]

      const props = mapStateToProps(state)
      expect(props.sites).toEqual(state.site.items)
    })
  })

  describe('mapDispatchToProps', () => {
    let dispatch

    beforeEach(() => {
      dispatch = jest.fn()
      actions.getAllSites = jest.fn().mockReturnValue({
        type: 'GET_ALL_SITES'
      })
    })

    test('maps store dispatch and actions to component props', () => {
      const props = mapDispatchToProps(dispatch)
      props.actions.getAllSites()

      expect(actions.getAllSites).toHaveBeenCalled()
      expect(dispatch).toHaveBeenCalledWith({
        type: 'GET_ALL_SITES'
      })
    })
  })
})
