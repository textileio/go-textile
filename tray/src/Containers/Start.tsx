import React, { Component } from 'react'
import { Button, Header, Image, Segment } from 'semantic-ui-react'
import { RouteComponentProps } from '@reach/router'
import SecurityImage from '../Assets/permissions@3x.png'
import BackArrow from '../Components/BackArrow'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

@connect('store') @observer
export default class Onboarding extends ConnectedComponent<RouteComponentProps, Stores> {
  handleCreate = () => this.props.navigate && this.props.navigate('./create')
  handleLogin = () => this.props.navigate && this.props.navigate('./login')
  render() {
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic>
          <Image centered size='medium' src={SecurityImage} />
          <Header as='h3'>
            Textile is a digital wallet for your personal data.
            <Header.Subheader>
              Textile helps you safely and securly store your files, chats, photos, and more.
              Putting control of your data back where it belongs... with you.
            </Header.Subheader>
          </Header>
        </Segment>
        <Button.Group fluid widths='2' style={{ position: 'absolute', bottom: 0 }}>
          <Button style={{ borderRadius: 0 }} positive content='Create' icon='key' type='button' onClick={this.handleCreate} />
          <Button style={{ borderRadius: 0 }} content='Sign-in' icon='sign-in' type='button' onClick={this.handleLogin} />
        </Button.Group>
        <BackArrow name='close' onClick={() => { this.stores.store.sendMessage({ name: 'quit' }) }} />
      </div>
    )
  }
}