import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import SiteList from '../SiteList'

describe('SiteList', () => {
  const props = {
    sites: []
  }

  test('renders react component with props', () => {
    props.sites = [
      { id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] },
      { id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }
    ]
    const wrapper = shallow(<SiteList {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('thead tr th').at(0).text()).toBe('Name')
    expect(wrapper.find('thead tr th').at(1).text()).toBe('Location')
    expect(wrapper.find('thead tr th').at(2).text()).toBe('Priority')
    expect(wrapper.find('thead tr th').at(3).text()).toBe('Tags')
    expect(wrapper.find('SiteListItem').get(0).props.site).toEqual(props.sites[0])
    expect(wrapper.find('SiteListItem').get(1).props.site).toEqual(props.sites[1])
  })
  test('renders react component with props', () => {
    props.sites = [
      { id: '3333-3333', name: 'Station 3', location: 'Ottawa, ON', priority: 4, tags: [ 'helium', 'nitrogen' ] },
      { id: '4444-4444', name: 'Station 4', location: 'Montreal, QC', priority: 5, tags: [ 'carbon', 'iron' ] }
    ]
    const wrapper = shallow(<SiteList {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find('thead tr th').at(0).text()).toBe('Name')
    expect(wrapper.find('thead tr th').at(1).text()).toBe('Location')
    expect(wrapper.find('thead tr th').at(2).text()).toBe('Priority')
    expect(wrapper.find('thead tr th').at(3).text()).toBe('Tags')
    expect(wrapper.find('SiteListItem').get(0).props.site).toEqual(props.sites[0])
    expect(wrapper.find('SiteListItem').get(1).props.site).toEqual(props.sites[1])
  })
})
