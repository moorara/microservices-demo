const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')
const MiniCssExtractPlugin = require("mini-css-extract-plugin")
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin")
const CleanWebpackPlugin = require('clean-webpack-plugin')
const InterpolateHtmlPlugin = require('./config/webpack/interpolate-html-plugin')

module.exports = [
  {
    name: 'client',
    target: 'web',
    mode: 'production',
    devtool: 'source-map',

    context: path.resolve(__dirname),
    entry: {
      app: './src/main.js'
    },

    resolve: {
      modules: [ 'node_modules' ],
      extensions: [ '.js', '.json', '.jsx' ],
      alias: {}
    },

    output: {
      publicPath: 'static/',
      path: path.resolve(__dirname, 'public'),
      filename: '[name].[chunkhash:8].bundle.js',
      chunkFilename: '[name].[chunkhash:8].chunk.js'
    },

    performance: {
      hints: 'warning',
      maxAssetSize: 500000
    },

    optimization: {
      minimizer: [
        new UglifyJsPlugin({
          cache: true,
          parallel: true,
          sourceMap: true
        }),
        new OptimizeCSSAssetsPlugin({})
      ],
      splitChunks: {
        cacheGroups: {
          vendors: {
            name: 'vendors',
            test: /[\\/]node_modules[\\/]/,
            chunks: 'all'
          },
          styles: {
            name: 'styles',
            test: /\.css$/,
            chunks: 'all',
            enforce: true
          }
        }
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
          use: [ MiniCssExtractPlugin.loader, 'css-loader' ]
        },
        {
          test: /\.less$/,
          use: [ MiniCssExtractPlugin.loader, 'css-loader', 'less-loader' ]
        },
        {
          test: /\.(scss|sass)$/,
          use: [ MiniCssExtractPlugin.loader, 'css-loader', 'sass-loader' ]
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
      new InterpolateHtmlPlugin({
        PAGE_TITLE: 'Control Center'
      }),
      new webpack.DefinePlugin({
        'process.env': {
          'NODE_ENV': JSON.stringify('production')
        }
      }),
      new MiniCssExtractPlugin({
        filename: "[name].[contenthash:8].css"
      }),
      new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/) // Required if using Moment.js
    ]
  }
]
