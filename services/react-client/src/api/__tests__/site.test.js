// This is imported from mock implementation
import * as axios from 'axios'

jest.mock('axios')

describe('siteApi', () => {
  beforeEach(() => {
    axios._clear()
  })

  afterEach(() => {
    expect(axios._create).toHaveBeenCalledWith({
      baseURL: '/api/v1/'
    })
  })

  describe('create', () => {
    let site

    beforeEach(() => {
      site = { name: 'Power Plant', location: 'Ottawa, ON', tags: [ 'hydro' ], priority: 3 }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.create(site).catch(e => {
        expect(e).toEqual(err)
        expect(axios._post).toHaveBeenCalledWith('sites', site)
      })
    })
    test('resolves with created site when api call succeeds', () => {
      const newSite = Object.assign({}, site, { id: '1111' })
      axios._mock(null, { status: 201, data: newSite })
      const siteApi = require('../site').default
      return siteApi.create(site).then(result => {
        expect(result).toEqual(newSite)
        expect(axios._post).toHaveBeenCalledWith('sites', site)
      })
    })
  })

  describe('all', () => {
    let sites

    beforeEach(() => {
      sites = [
        { id: '1111', name: 'Power Plant', location: 'Ottawa, ON', tags: [ 'hydro' ], priority: 3 },
        { id: '2222', name: 'Gas Station', location: 'Ottawa, ON', tags: [ 'fuel' ], priority: 2 },
      ]
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.all().catch(e => {
        expect(e).toEqual(err)
        expect(axios._get).toHaveBeenCalledWith('sites')
      })
    })
    test('resolves with all sites when api call succeeds', () => {
      axios._mock(null, { status: 200, data: sites })
      const siteApi = require('../site').default
      return siteApi.all().then(result => {
        expect(result).toEqual(sites)
        expect(axios._get).toHaveBeenCalledWith('sites')
      })
    })
  })

  describe('get', () => {
    let id, site

    beforeEach(() => {
      id = '1111'
      site = { id, name: 'Power Plant', location: 'Ottawa, ON', tags: [ 'hydro' ], priority: 3 }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.get(id).catch(e => {
        expect(e).toEqual(err)
        expect(axios._get).toHaveBeenCalledWith(`sites/${id}`)
      })
    })
    test('resolves with requested site when api call succeeds', () => {
      axios._mock(null, { status: 200, data: site })
      const siteApi = require('../site').default
      return siteApi.get(id).then(result => {
        expect(result).toEqual(site)
        expect(axios._get).toHaveBeenCalledWith(`sites/${id}`)
      })
    })
  })

  describe('update', () => {
    let site

    beforeEach(() => {
      site = { id: '1111', name: 'Power Station', location: 'Ottawa, ON, CANADA', tags: [ 'power', 'hydro' ], priority: 4 }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.update(site).catch(e => {
        expect(e).toEqual(err)
        expect(axios._put).toHaveBeenCalledWith(`sites/${site.id}`, site)
      })
    })
    test('resolves successfully when api call succeeds', () => {
      axios._mock(null, { status: 204, data: null })
      const siteApi = require('../site').default
      return siteApi.update(site).then(result => {
        expect(result).toBeNull()
        expect(axios._put).toHaveBeenCalledWith(`sites/${site.id}`, site)
      })
    })
  })

  describe('modify', () => {
    let site

    beforeEach(() => {
      site = { id: '1111', name: 'Power Station', location: 'Ottawa, ON, CANADA' }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.modify(site).catch(e => {
        expect(e).toEqual(err)
        expect(axios._patch).toHaveBeenCalledWith(`sites/${site.id}`, site)
      })
    })
    test('resolves with modified site when api call succeeds', () => {
      const modifiedSite = { id: site.id, name: 'Power Station', location: 'Ottawa, ON, CANADA', tags: [ 'hydro' ], priority: 3 }
      axios._mock(null, { status: 200, data: modifiedSite })
      const siteApi = require('../site').default
      return siteApi.modify(site).then(result => {
        expect(result).toEqual(modifiedSite)
        expect(axios._patch).toHaveBeenCalledWith(`sites/${site.id}`, site)
      })
    })
  })

  describe('delete', () => {
    let id

    beforeEach(() => {
      id = '1111'
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const siteApi = require('../site').default
      return siteApi.delete(id).catch(e => {
        expect(e).toEqual(err)
        expect(axios._delete).toHaveBeenCalledWith(`sites/${id}`)
      })
    })
    test('resolves successfully when api call succeeds', () => {
      axios._mock(null, { status: 204, data: null })
      const siteApi = require('../site').default
      return siteApi.delete(id).then(result => {
        expect(result).toBeNull()
        expect(axios._delete).toHaveBeenCalledWith(`sites/${id}`)
      })
    })
  })
})
