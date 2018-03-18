const escapeStringRegexp = require('escape-string-regexp')

class InterpolateHtmlPlugin {
  constructor (replacements) {
    this.replacements = replacements
  }

  apply (compiler) {
    compiler.hooks.compilation.tap('InterpolateHtmlPlugin', compilation => {
      compilation.hooks.htmlWebpackPluginBeforeHtmlProcessing.tap('InterpolateHtmlPlugin', data => {
        Object.keys(this.replacements).forEach(key => {
          const value = this.replacements[key]
          data.html = data.html.replace(
            new RegExp(`%${escapeStringRegexp(key)}%`, 'g'),
            value
          )
        })
      })
    })
  }
}

module.exports = InterpolateHtmlPlugin
