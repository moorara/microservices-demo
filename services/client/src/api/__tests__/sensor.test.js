// This is imported from mock implementation
import * as axios from 'axios'

jest.mock('axios')

describe('sensorApi', () => {
  beforeEach(() => {
    axios._clear()
  })

  afterEach(() => {
    expect(axios._create).toHaveBeenCalledWith({
      baseURL: '/api/v1/'
    })
  })

  describe('create', () => {
    let sensor

    beforeEach(() => {
      sensor = { siteId: '1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
    })

    test('should reject with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const sensorApi = require('../sensor').default
      return sensorApi.create(sensor).catch(e => {
        expect(e).toEqual(err)
        expect(axios._post).toHaveBeenCalledWith('sensors', sensor)
      })
    })
    test('resolves with created sensor when api call succeeds', () => {
      const newSensor = Object.assign({}, sensor, { id: 'aaaa' })
      axios._mock(null, { status: 201, data: newSensor })
      const sensorApi = require('../sensor').default
      return sensorApi.create(sensor).then(result => {
        expect(result).toEqual(newSensor)
        expect(axios._post).toHaveBeenCalledWith('sensors', sensor)
      })
    })
  })

  describe('all', () => {
    let siteId, sensors

    beforeEach(() => {
      siteId = '1111'
      sensors = [
        { id: 'aaaa', siteId, name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 },
        { id: 'bbbb', siteId, name: 'temperature', unit: 'farenheit', minSafe: -22, maxSafe: 86 }
      ]
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const sensorApi = require('../sensor').default
      return sensorApi.all(siteId).catch(e => {
        expect(e).toEqual(err)
        expect(axios._get).toHaveBeenCalledWith(`sensors?siteId=${siteId}`)
      })
    })
    test('resolves with all sensors when api call succeeds', () => {
      axios._mock(null, { status: 200, data: sensors })
      const sensorApi = require('../sensor').default
      return sensorApi.all(siteId).then(result => {
        expect(result).toEqual(sensors)
        expect(axios._get).toHaveBeenCalledWith(`sensors?siteId=${siteId}`)
      })
    })
  })

  describe('get', () => {
    let id, sensor

    beforeEach(() => {
      id = 'aaaa'
      sensor = { id, siteId: '1111', name: 'temperature', unit: 'celsius', minSafe: -30, maxSafe: 30 }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const sensorApi = require('../sensor').default
      return sensorApi.get(id).catch(e => {
        expect(e).toEqual(err)
        expect(axios._get).toHaveBeenCalledWith(`sensors/${id}`)
      })
    })
    test('resolves with requested sensor when api call succeeds', () => {
      axios._mock(null, { status: 200, data: sensor })
      const sensorApi = require('../sensor').default
      return sensorApi.get(id).then(result => {
        expect(result).toEqual(sensor)
        expect(axios._get).toHaveBeenCalledWith(`sensors/${id}`)
      })
    })
  })

  describe('update', () => {
    let sensor

    beforeEach(() => {
      sensor = { id: 'aaaa', siteId: '1111', name: 'temperature', unit: 'farenheit', minSafe: -22, maxSafe: 86 }
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const sensorApi = require('../sensor').default
      return sensorApi.update(sensor).catch(e => {
        expect(e).toEqual(err)
        expect(axios._put).toHaveBeenCalledWith(`sensors/${sensor.id}`, sensor)
      })
    })
    test('resolves successfully when api call succeeds', () => {
      axios._mock(null, { status: 204, data: null })
      const sensorApi = require('../sensor').default
      return sensorApi.update(sensor).then(result => {
        expect(result).toBeNull()
        expect(axios._put).toHaveBeenCalledWith(`sensors/${sensor.id}`, sensor)
      })
    })
  })

  describe('delete', () => {
    let id

    beforeEach(() => {
      id = 'aaaa'
    })

    test('rejects with error when api call fails', () => {
      const err = new Error('error')
      axios._mock(err)
      const sensorApi = require('../sensor').default
      return sensorApi.delete(id).catch(e => {
        expect(e).toEqual(err)
        expect(axios._delete).toHaveBeenCalledWith(`sensors/${id}`)
      })
    })
    test('resolves successfully when api call succeeds', () => {
      axios._mock(null, { status: 204, data: null })
      const sensorApi = require('../sensor').default
      return sensorApi.delete(id).then(result => {
        expect(result).toBeNull()
        expect(axios._delete).toHaveBeenCalledWith(`sensors/${id}`)
      })
    })
  })
})
