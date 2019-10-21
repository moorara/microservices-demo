const express = require('express')

const router = express.Router()

router.use('/readiness', (req, res) => {
  res.sendStatus(200)
})

module.exports = router
