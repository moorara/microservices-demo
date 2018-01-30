const _ = require('lodash')
const mongoose = require('mongoose')

const MODEL_NAME = 'Link'
const PUBLIC_PROPERTIES = [ 'id', 'url', 'title', 'tags', 'rank' ]

class Link {
  constructor (config, options) {
    options = options || {}
    this.Model = this._model()
  }

  _model () {
    if (mongoose.models[MODEL_NAME]) {
      return mongoose.model(MODEL_NAME)
    }

    const schema = new mongoose.Schema({
      url: { type: String, required: true, unique: true, index: true },
      title: { type: String, required: true },
      tags: [ String ],
      rank: Number,
      createdAt: { type: Date, default: Date.now },
      updatedAt: { type: Date, default: Date.now }
    })

    schema.set('toJSON', {
      transform: doc => _.pick(doc, PUBLIC_PROPERTIES)
    })

    schema.pre('save', function (next) {
      this.updatedAt = Date.now()
      next()
    })

    schema.methods.formatTitle = function () {
      return this.title.replace(/\b\w/g, l => l.toUpperCase())
    }

    return mongoose.model(MODEL_NAME, schema)
  }
}

module.exports = Link
