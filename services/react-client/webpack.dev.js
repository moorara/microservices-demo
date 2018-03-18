// https://webpack.js.org/concepts
// https://webpack.js.org/configuration
// https://webpack.js.org/loaders
// https://webpack.js.org/plugins

const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const InterpolateHtmlPlugin = require('./config/webpack/interpolate-html-plugin')

const port = process.env.PORT || 4000
const apiPort = process.env.API_PORT || 4001

module.exports = {
  target: 'web',
  mode: 'development',
  devtool: 'eval-source-map',

  entry: {
    app: path.resolve(__dirname, 'src/main.js')
  },

  output: {
    path: path.resolve(__dirname, 'public'),
    filename: '[name].bundle.js',
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
        use: [ 'style-loader', 'css-loader' ]
      },
      {
        test: /\.less$/,
        use: [ 'style-loader', 'css-loader', 'less-loader' ]
      },
      {
        test: /\.(scss|sass)$/,
        use: [ 'style-loader', 'css-loader', 'sass-loader' ]
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

  performance: {
    hints: false
  },

  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'src/index.html'),
      filename: 'index.html',
      inject: 'body',
      favicon: 'src/favicon.ico'
    }),
    new InterpolateHtmlPlugin({
      PAGE_TITLE: 'Control Center (dev)'
    }),
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify('development')
      }
    }),
    new webpack.NamedModulesPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/) // Required if using Moment.js
  ],

  devServer: {
    contentBase: path.resolve(__dirname, 'public'),
    publicPath: '/',
    host: 'localhost',
    port: port,
    compress: true,
    hot: true,
    noInfo: true,
    open: true,
    historyApiFallback: true,
    overlay: {
      errors: true,
      warnings: true
    },
    proxy: {
      '/api/v1': {
        target: `http://localhost:${apiPort}`,
        pathRewrite: { '^/api/v1': '' }
      }
    }
  }
}
