import React from 'react'
import PropTypes from 'prop-types'

const TextInput = ({ name, label, placeholder, value, isOk, isError, help, onChange }) => {
  let inputClass = 'input'
  let helpClass = 'help'

  if (isOk) {
    inputClass += ' is-success'
    helpClass += ' is-success'
  } else if (isError) {
    inputClass += ' is-danger'
    helpClass += ' is-danger'
  }

  return (
    <div className="field">
      <label className="label" htmlFor={name}>{label}</label>
      <div className="control has-icons-right">
        <input
          type="text"
          className={inputClass}
          name={name}
          placeholder={placeholder}
          value={value}
          onChange={onChange} />
        { isOk &&
          <span className="icon is-small is-right">
            <i className="fa fa-check-circle" />
          </span>
        }
        { isError &&
          <span className="icon is-small is-right">
            <i className="fa fa-exclamation-circle" />
          </span>
        }
      </div>
      { help && <p className={helpClass}>{help}</p> }
    </div>
  )
}

TextInput.propTypes = {
  name: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  placeholder: PropTypes.string,
  value: PropTypes.string,
  isOk: PropTypes.bool,
  isError: PropTypes.bool,
  help: PropTypes.string,
  onChange: PropTypes.func
}

TextInput.defaultProps = {
  placeholder: '',
  value: '',
  isOk: false,
  isError: false,
  help: ''
}

export default TextInput
