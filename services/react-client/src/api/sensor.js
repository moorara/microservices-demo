import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1/'
})

export default {
  async create (sensor) {
    try {
      const res = await http.post('sensors', sensor)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async all (siteId) {
    try {
      const res = await http.get(`sensors?siteId=${siteId}`)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async get (id) {
    try {
      const res = await http.get(`sensors/${id}`)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async update (sensor) {
    try {
      const res = await http.put(`sensors/${sensor.id}`, sensor)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async delete (id) {
    try {
      const res = await http.delete(`sensors/${id}`)
      return res.data
    } catch (err) {
      throw err
    }
  },
}
