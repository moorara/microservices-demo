import React from 'react'
import PropTypes from 'prop-types'

class Loading extends React.Component {
  static defaultProps = {
    dots: 5,
    interval: 200
  }

  static propTypes = {
    dots: PropTypes.number,
    interval: PropTypes.number
  }

  constructor (props, context) {
    super(props, context)
    this.state = {
      frame: 1
    }
  }

  componentDidMount () {
    this.interval = setInterval(() => {
      this.setState({
        frame: this.state.frame + 1
      })
    }, this.props.interval)
  }

  componentWillUnmount () {
    clearInterval(this.interval)
  }

  render () {
    let text = ''
    for (
      let dots = this.state.frame % (this.props.dots + 1);
      dots > 0;
      dots--, text += '•'
    );

    return (
      <div className="dots">
        <span>{text}</span>
      </div>
    )
  }
}

export default Loading
