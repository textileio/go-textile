import React from 'react'
import { RouteComponentProps } from '@reach/router'
import BackArrow from '../Components/BackArrow'
import { Image, Dimmer } from 'semantic-ui-react'
import Pulse from 'react-reveal/Pulse'
import Logo from '../Assets/LaunchLogo@3x.png'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

@connect('store') @observer
export default class Splash extends ConnectedComponent<RouteComponentProps, Stores> {
  render() {
    return (
      <Dimmer inverted active>
        <Pulse forever>
          <Image centered verticalAlign='middle' size='small' src={Logo} />
        </Pulse>
        <BackArrow name='close' onClick={() => { this.stores.store.sendMessage({ name: 'quit' }) }} />
      </Dimmer>
    )
  }
}
