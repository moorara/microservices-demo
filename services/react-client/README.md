# react-client

## Commands

| Command                  | Description                                                   |
|--------------------------|---------------------------------------------------------------|
| `yarn run lint`          | Runs linter with configured rules                             |
| `yarn run lint:watch`    | Runs linter and watches for changes                           |
| `yarn run test`          | Runs unit tests                                               |
| `yarn run test:update`   | Runs unit tests and updates `jest snapshots`                  |
| `yarn run test:watch`    | Runs unit tests and watches for changes                       |
| `yarn run test:coverage` | Runs unit tests and generate coverage reports                 |
| `yarn run dev:node`      | Starts development server using `node.js`                     |
| `yarn run dev:api`       | Starts development API server using `json-server`             |
| `yarn run dev:webpack`   | Starts development server using built-in `webpack dev server` |
| `yarn run build:webpack` | Builds application for production using `webpack cli`         |
| `yarn run build:node`    | Builds application for production using `node.js`             |
| `yarn run prod:server`   | Runs a production server for serving the application          |
| `make docker`            | Builds Docker image                                           |
| `make up`                | Runs the service locally in containers                        |
| `make down`              | Stops and removes local containers                            |

## Guides

### Snapshot Testing

When it comes to snapshot testing React components with `jest`, there are two options:

  1. Snapshot testing with `react-test-renderer`
  2. Snapshot testing with `enzyme` and `enzyme-to-json`

### Resources

  * **Webpack**
    - https://webpack.js.org/concepts
    - https://webpack.js.org/configuration
    - https://webpack.js.org/guides
    - https://webpack.js.org/loaders
    - https://webpack.js.org/plugins
    - https://github.com/jantimon/html-webpack-plugin
    - https://github.com/kangax/html-minifier
    - https://github.com/webpack/webpack-dev-middleware
    - https://github.com/glenjamin/webpack-hot-middleware

  * **Babel**
    - https://babeljs.io/docs/plugins
    - https://babeljs.io/docs/plugins/preset-env
    - https://babeljs.io/docs/plugins/preset-react
    - https://babeljs.io/docs/plugins/preset-stage-2
    - https://babeljs.io/docs/plugins/transform-runtime

  * **ESLint**
    - https://eslint.org/docs/user-guide/getting-started
    - https://eslint.org/docs/user-guide/configuring
    - https://eslint.org/docs/rules
    - https://www.npmjs.com/package/eslint-plugin-react
    - https://www.npmjs.com/package/eslint-config-standard
    - https://www.npmjs.com/package/eslint-config-standard-react

  * **React**
    - https://reactjs.org/docs/forms.html
    - https://reactjs.org/docs/fragments.html

  * **Redux**
    - https://redux.js.org/advanced/async-actions
    - https://redux.js.org/recipes/structuring-reducers/initializing-state
    - https://redux.js.org/recipes/server-rendering
    - https://redux.js.org/recipes/writing-tests

  * **Jest**
    - https://facebook.github.io/jest/docs/en/jest-object.html
    - https://facebook.github.io/jest/docs/en/mock-function-api.html
    - https://facebook.github.io/jest/docs/en/mock-functions.html
    - https://facebook.github.io/jest/docs/en/tutorial-async.html
    - https://facebook.github.io/jest/docs/en/manual-mocks.html
    - https://facebook.github.io/jest/docs/en/es6-class-mocks.html
    - https://facebook.github.io/jest/docs/en/snapshot-testing.html
    - https://facebook.github.io/jest/docs/en/configuration.html

  * **Enzyme**
    - http://airbnb.io/enzyme/docs/api/shallow.html
    - http://airbnb.io/enzyme/docs/api/mount.html
    - http://airbnb.io/enzyme/docs/api/render.html
    - http://airbnb.io/enzyme/docs/api/selector.html
    - https://github.com/adriantoine/enzyme-to-json
    - https://github.com/airbnb/enzyme/tree/master/packages/enzyme-adapter-react-16

  * **Misc**
    - https://github.com/typicode/json-server
    - https://www.jstwister.com/post/jest-snapshot-testing-with-enzyme
