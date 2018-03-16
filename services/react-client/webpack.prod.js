// https://webpack.js.org/concepts
// https://webpack.js.org/configuration
// https://webpack.js.org/loaders
// https://webpack.js.org/plugins

const path = require('path')
const webpack = require('webpack')
const CleanWebpackPlugin = require('clean-webpack-plugin')
const InterpolateHtmlPlugin = require('interpolate-html-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')

const extractCSS = new ExtractTextPlugin('[name].[contenthash:8].css')
const extractLESS = new ExtractTextPlugin('[name].[contenthash:8].css')
const extractSASS = new ExtractTextPlugin('[name].[contenthash:8].css')

module.exports = {
  target: 'web',
  mode: 'production',
  devtool: 'source-map',

  entry: {
    vendor: [
      'axios', 'bulma',
      'react', 'react-dom', 'react-router-dom',
      'redux', 'react-redux', 'redux-thunk', 'redux-logger',
    ],
    app: path.resolve(__dirname, 'src/main.js')
  },

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: '[name].[hash:8].bundle.js',
    publicPath: '/'
  },

  resolve: {
    modules: [ 'node_modules' ],
    extensions: [ '.js', '.json', '.jsx' ],
    alias: {
    }
  },

  module: {
    rules: [
      // JavaScript
      {
        test: /\.(js|jsx)$/,
        include: path.resolve(__dirname, 'src'),
        exclude: path.resolve(__dirname, 'node_modules'),
        use: [ 'babel-loader', 'eslint-loader' ]
      },

      // Stylesheets
      {
        test: /\.css$/,
        use: extractCSS.extract({
          fallback: 'style-loader',
          use: [ 'css-loader' ]
        })
      },
      {
        test: /\.less$/,
        use: extractLESS.extract({
          fallback: 'style-loader',
          use: [ 'css-loader', 'less-loader' ]
        })
      },
      {
        test: /\.(scss|sass)$/,
        use: extractSASS.extract({
          fallback: 'style-loader',
          use: [ 'css-loader', 'sass-loader' ]
        })
      },

      // Fonts
      {
        test: /\.(ttf|otf|eot|woff2?)(\?.*)?$/,
        use: [{
          loader: 'url-loader',
          options: {
            limit: 4096,
            fallback: 'file-loader',
            name: 'fonts/[name].[hash:8].[ext]' // option for file-loader fallback
          }
        }]
      },

      // Images
      {
        test: /\.(png|jpg|jpeg|gif|svg)(\?.*)?$/,
        use: [{
          loader: 'url-loader',
          options: {
            limit: 4096,
            fallback: 'file-loader',
            name: 'images/[name].[hash:8].[ext]' // option for file-loader fallback
          }
        }]
      },

      // Media
      {
        test: /\.(mp3|wav|aac|flac|mp4|ogg|webm)(\?.*)?$/,
        use: [{
          loader: 'file-loader',
          options: {
            name: 'media/[name].[hash:8].[ext]'
          }
        }]
      }
    ]
  },

  plugins: [
    new CleanWebpackPlugin([ 'public' ]),
    new InterpolateHtmlPlugin({
      'PAGE_TITLE': 'Control Center'
    }),
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'src/index.html'),
      filename: path.resolve(__dirname, 'public/index.html'),
      inject: 'body',
      favicon: 'src/favicon.ico',
      minify: {
        html5: true,
        minifyJS: true,
        minifyCSS: true,
        minifyURLs: true,
        keepClosingSlash: true,
        collapseWhitespace: true,
        removeComments: true,
        removeEmptyAttributes: true,
        removeRedundantAttributes: true
      }
    }),
    new webpack.optimize.CommonsChunkPlugin({
      names: [ 'vendor' ]
    }),
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify('production')
      }
    }),
    new UglifyJsPlugin({
      sourceMap: true,
      uglifyOptions: {
        output: {
          ascii_only: true
        },
        compress: {
          comparisons: false
        },
        mangle: {
          safari10: true
        }
      }
    }),
    extractCSS, extractLESS, extractSASS,
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/) // Required if using Moment.js
  ]
}
