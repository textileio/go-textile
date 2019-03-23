import React, { Component } from 'react'
import { Image } from 'semantic-ui-react'
import Logo from './LaunchLogo@3x.png'
import Pulse from 'react-reveal/Pulse'

class Splash extends Component {
  render() {
    return (
      <Pulse forever>
        <Image style={{ margin: 'auto' }} centered size='medium' src={Logo} />
      </Pulse>
    )
  }
}

export default Splash
