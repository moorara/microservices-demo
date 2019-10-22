# react-client

## Commands

| Command                  | Description                                               |
|--------------------------|-----------------------------------------------------------|
| `yarn run lint`          | Run linter with configured rules                          |
| `yarn run lint:fix`      | Fix linting issues                                        |
| `yarn run lint:watch`    | Run linter and watches for changes                        |
| `yarn run test`          | Run unit tests in the interactive watch mode              |
| `yarn run test:coverage` | Run unit tests and generate coverage reports              |
| `yarn run test:update`   | Run unit tests and updates `jest snapshots`               |
| `yarn run start`         | Start serving application in development mode             |
| `yarn run start:proxy`   | Start a proxy server for serving application and mock API |
| `yarn run start:api`     | Start a mock server for serving API requests              |
| `yarn run build`         | Build application bundle for production                   |
| `make docker`            | Build Docker image                                        |
| `make up`                | Run the service locally in containers                     |
| `make down`              | Stop and removes local containers                         |
| `make push`              | Push built images to registry                             |
| `make save-images`       | Save built images to disk                                 |

## Documentation

  - https://webpack.js.org
  - https://babeljs.io
  - https://eslint.org
  - https://reactjs.org
    - https://facebook.github.io/create-react-app
      - https://facebook.github.io/create-react-app/docs/proxying-api-requests-in-development
  - https://redux.js.org
  - https://jestjs.io
  - http://airbnb.io/enzyme

### Snapshot Testing

When it comes to snapshot testing React components with `jest`, there are two options:

  1. Snapshot testing with `react-test-renderer`
  2. Snapshot testing with `enzyme` and `enzyme-to-json`
