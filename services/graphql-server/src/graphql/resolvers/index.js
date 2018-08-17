const _ = require('lodash')

const siteResolvers = require('./site')
const sensorResolvers = require('./sensor')
const switchResolvers = require('./switch')

module.exports = _.merge({}, siteResolvers, sensorResolvers, switchResolvers)
