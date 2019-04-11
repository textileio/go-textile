import React from 'react'
import { Dimmer, Loader } from 'semantic-ui-react'
import { RouteComponentProps } from '@reach/router'
import BackArrow from '../Components/BackArrow'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

@connect('store') @observer
export default class Loading extends ConnectedComponent<RouteComponentProps, Stores> {
  render() {
    return (
      <Dimmer inverted active>
        <Loader size='massive' />
        <BackArrow name='close' onClick={() => { this.stores.store.sendMessage({ name: 'quit' }) }} />
      </Dimmer>
    )
  }
}
