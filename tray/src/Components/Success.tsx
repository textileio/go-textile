import React from 'react'
import { Button, Header, Segment, Image } from 'semantic-ui-react'
// import { Fade } from 'react-reveal'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'
import { RouteComponentProps } from '@reach/router'
import AccountsImage from '../Assets/new_avatar@3x.png'

@connect('store') @observer
export default class Onboarding extends ConnectedComponent<RouteComponentProps, Stores> {
  handleClose = () => this.stores.store.sendMessage({ name: 'hide' })
  handleLink = () => this.stores.store.sendMessage({ name: 'open', payload: { url: 'http://127.0.0.1:40602/docs/index.html' } })
  render() {
    const { store } = this.stores
    return (
      // <Fade>
        <div>
          <Segment basic attached>
            <Image centered size='small' src={AccountsImage} />
            <Header as='h3'>
              Success! You're all set to go.
            </Header>
            <p>
              You can close this window and you start interacting with the tray app to create
              new accounts, turn your node on/off, and edit your node's config file.... blah
            </p>
            <p>
              Here's your account address:
              {store.profile.address} <a onClick={this.handleLink}>and here are the API docs</a>.
            </p>
          </Segment>
          <Button.Group attached='bottom'>
            <Button content='Close' icon='window close' type='button' onClick={this.handleClose} />
          </Button.Group>
        </div>
      // </Fade>
    )
  }
}