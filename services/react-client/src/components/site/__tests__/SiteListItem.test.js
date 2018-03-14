import React from 'react'
import { shallow } from 'enzyme'
import { shallowToJson } from 'enzyme-to-json'

import SiteListItem from '../SiteListItem'

describe('SiteListItem', () => {
  const props = {
    site: {}
  }

  test('renders react component with props', () => {
    props.site = { id: '1111-1111', name: 'Station 1', location: 'Toronto, ON', priority: 2, tags: [ 'oxygen', 'hydrogen' ] }
    const wrapper = shallow(<SiteListItem {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find(`Link [to="/sites/${props.site.id}"]`))
    expect(wrapper.find('td').at(1).text()).toBe(`${props.site.location}`)
    expect(wrapper.find('td').at(2).text()).toBe(`${props.site.priority}`)
    expect(wrapper.find('td').at(3).text()).toBe(`${props.site.tags.join(' ')}`)
  })
  test('renders react component with props', () => {
    props.site = { id: '2222-2222', name: 'Station 2', location: 'Kingston, ON', priority: 3, tags: [ 'gold', 'silver' ] }
    const wrapper = shallow(<SiteListItem {...props} />)

    expect(shallowToJson(wrapper)).toMatchSnapshot()
    expect(wrapper.find(`Link [to="/sites/${props.site.id}"]`))
    expect(wrapper.find('td').at(1).text()).toBe(`${props.site.location}`)
    expect(wrapper.find('td').at(2).text()).toBe(`${props.site.priority}`)
    expect(wrapper.find('td').at(3).text()).toBe(`${props.site.tags.join(' ')}`)
  })
})
