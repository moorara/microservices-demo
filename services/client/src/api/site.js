import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1/'
})

export default {
  async create (site) {
    try {
      const res = await http.post('sites', site)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async all () {
    try {
      const res = await http.get('sites')
      return res.data
    } catch (err) {
      throw err
    }
  },

  async get (id) {
    try {
      const res = await http.get(`sites/${id}`)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async update (site) {
    try {
      const res = await http.put(`sites/${site.id}`, site)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async modify (site) {
    try {
      const res = await http.patch(`sites/${site.id}`, site)
      return res.data
    } catch (err) {
      throw err
    }
  },

  async delete (id) {
    try {
      const res = await http.delete(`sites/${id}`)
      return res.data
    } catch (err) {
      throw err
    }
  },
}
