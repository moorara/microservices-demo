/*
 * This is for building the client application bundle using Node.js
 */

import chalk from 'chalk'
import webpack from 'webpack'

import webpackConfig from '../webpack.prod'

webpack(webpackConfig, (err, stats) => {
  if (err) {
    console.log(chalk.red.bold('Build failed with errors.\n'))
    throw err
  }

  process.stdout.write(stats.toString({
    colors: true
  }) + '\n\n')

  if (stats.hasErrors()) {
    console.log(chalk.red.bold('Build failed with errors.\n'))
    process.exit(1)
  }

  console.log(chalk.green.bold('Build complete successfully.\n'))
})
