let mockErr = null
let mockRes = { status: 200, data: null }

export const _create = jest.fn().mockImplementation(function () {
  return this
})

export const _get = jest.fn().mockImplementation(() => {
  if (mockErr) {
    return Promise.reject(mockErr)
  }
  return Promise.resolve(mockRes)
})

export const _post = jest.fn().mockImplementation(() => {
  if (mockErr) {
    return Promise.reject(mockErr)
  }
  return Promise.resolve(mockRes)
})

export const _put = jest.fn().mockImplementation(() => {
  if (mockErr) {
    return Promise.reject(mockErr)
  }
  return Promise.resolve(mockRes)
})

export const _patch = jest.fn().mockImplementation(() => {
  if (mockErr) {
    return Promise.reject(mockErr)
  }
  return Promise.resolve(mockRes)
})

export const _delete = jest.fn().mockImplementation(() => {
  if (mockErr) {
    return Promise.reject(mockErr)
  }
  return Promise.resolve(mockRes)
})

export const _mock = (err, res) => {
  mockErr = err
  mockRes = res
}

export const _clear = () => {
  _get.mockClear()
  _post.mockClear()
  _put.mockClear()
  _patch.mockClear()
  _delete.mockClear()
  mockErr = null
  mockRes = { status: 200, data: null }
}

export default {
  create: _create,
  get: _get,
  post: _post,
  put: _put,
  patch: _patch,
  delete: _delete
}
