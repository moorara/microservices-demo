const http = require('http')
const axios = require('axios')
const opentracing = require('opentracing')

const Logger = require('../utils/logger')
const { createTracer } = require('../utils/tracer')

const timeout = 1000

class SiteService {
  constructor (config, options) {
    options = options || {}
    this.logger = options.logger || new Logger('SiteService')
    this.histogram = options.histogram || { observe () {} }
    this.summary = options.summary || { observe () {} }
    this.tracer = options.tracer || createTracer({ serviceName: 'SiteService' })
    this.axios = options.axios || axios.create({
      timeout,
      httpAgent: new http.Agent({ keepAlive: true }),
      baseURL: `http://${config.siteServiceAddr}/v1/`
    })
  }

  async exec (context, name, func) {
    let result, err
    let latency

    // https://opentracing-javascript.surge.sh/interfaces/spanoptions.html
    // https://github.com/opentracing/specification/blob/master/semantic_conventions.md
    const span = this.tracer.startSpan(name, { childOf: context.span }) // { childOf: context.span.context() }
    span.setTag(opentracing.Tags.SPAN_KIND, 'client')
    span.setTag(opentracing.Tags.PEER_SERVICE, 'site-service')
    span.setTag(opentracing.Tags.PEER_ADDRESS, this.axios.defaults.baseUrl)

    const headers = {}
    this.tracer.inject(span, opentracing.FORMAT_HTTP_HEADERS, headers)

    // Core functionality
    try {
      const startTime = Date.now()
      result = await func(headers)
      latency = (Date.now() - startTime) / 1000
    } catch (e) {
      err = e
      this.logger.error(err)
    }

    // Metrics
    const labelValues = { op: name, success: err ? 'false' : 'true' }
    this.histogram.observe(labelValues, latency)
    this.summary.observe(labelValues, latency)

    // Tracing
    span.log({
      event: name,
      message: err ? err.message : 'successful!'
    })
    span.finish()

    if (err) {
      throw err
    }

    return result
  }

  create (context, input) {
    return this.exec(context, 'create-site', headers => {
      return this.axios.request({
        headers,
        method: 'post',
        url: '/sites',
        data: input
      }).then(res => res.data)
    })
  }

  all (context, query) {
    return this.exec(context, 'get-sites', headers => {
      return this.axios.request({
        headers,
        method: 'get',
        url: '/sites',
        params: query
      }).then(res => res.data)
    })
  }

  get (context, id) {
    return this.exec(context, 'get-site', headers => {
      return this.axios.request({
        headers,
        method: 'get',
        url: `/sites/${id}`
      }).then(res => res.data)
    })
  }

  update (context, id, input) {
    return this.exec(context, 'update-site', async headers => {
      try {
        await this.axios.request({
          headers,
          method: 'put',
          url: `/sites/${id}`,
          data: input
        })

        // site-service does not respond with updated site
        const updated = await this.axios.request({
          headers,
          method: 'get',
          url: `/sites/${id}`
        })

        return updated.data
      } catch (err) {
        throw err
      }
    })
  }

  modify (context, id, input) {
    return this.exec(context, 'modify-site', headers => {
      return this.axios.request({
        headers,
        method: 'patch',
        url: `/sites/${id}`,
        data: input
      }).then(res => res.data)
    })
  }

  delete (context, id) {
    return this.exec(context, 'delete-site', headers => {
      return this.axios.request({
        headers,
        method: 'delete',
        url: `/sites/${id}`
      }).then(res => res.data)
    })
  }
}

module.exports = SiteService
