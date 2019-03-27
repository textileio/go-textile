import React from 'react'
import { Image, Dimmer } from 'semantic-ui-react'
import Pulse from 'react-reveal/Pulse'
import Logo from '../Assets/LaunchLogo@3x.png'
import { RouteComponentProps } from '@reach/router'

const Splash = (props: RouteComponentProps) => {
  return (
    <Dimmer inverted active>
      <Pulse forever>
        <Image centered verticalAlign='middle' size='small' src={Logo} />
      </Pulse>
    </Dimmer>
  )
}

export default Splash
