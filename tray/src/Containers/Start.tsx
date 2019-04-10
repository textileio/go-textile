import React, { Component } from 'react'
import { Button, Header, Image, Segment } from 'semantic-ui-react'
import { RouteComponentProps } from '@reach/router'
import SecurityImage from '../Assets/permissions@3x.png'

export default class Onboarding extends Component<RouteComponentProps> {
  handleCreate = () => this.props.navigate && this.props.navigate('./create')
  handleLogin = () => this.props.navigate && this.props.navigate('./login')
  render() {
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic>
          <Image centered size='medium' src={SecurityImage} />
          <Header as='h3'>
            Textile is a digital wallet for your personal data.
        </Header>
          <p>
            Textile helps you safely and securly store your files, chats, photos, and more.
            Putting control of your data back where it belongs... with you.
        </p>
        </Segment>
        <Button.Group fluid style={{ position: 'absolute', bottom: 0 }}>
          <Button style={{ borderRadius: 0 }} positive content='Create' icon='key' type='button' onClick={this.handleCreate} />
          <Button style={{ borderRadius: 0 }} content='Sign-in' icon='sign-in' type='button' onClick={this.handleLogin} />
        </Button.Group>
      </div>
    )
  }
}