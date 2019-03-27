import React, { Component } from 'react'
import { Button, Header, Segment, Image } from 'semantic-ui-react'
import { Fade } from 'react-reveal'
import { RouteComponentProps } from '@reach/router'
import SecurityImage from '../Assets/permissions@3x.png'

export default class Onboarding extends Component<RouteComponentProps> {
  handleCreate = () => this.props.navigate && this.props.navigate('./create')
  handleLogin = () => this.props.navigate && this.props.navigate('./login')
  render() {
    return (
      <Fade>
        <Segment raised>
          <Segment basic attached>
            <Image centered size='medium' src={SecurityImage} />
            <Header as='h3'>
              Textile is a digital wallet for your personal data.
            </Header>
            <p>
              Textile helps you safely and securly store your files, chats, photos, and more.
              Putting control of your data back where it belongs... with you.
            </p>
          </Segment>
          <Button.Group attached='bottom'>
            <Button positive content='Create' icon='key' type='button' onClick={this.handleCreate} />
            <Button content='Sign-in' icon='sign-in' type='button' onClick={this.handleLogin} />
          </Button.Group>
        </Segment>
      </Fade>
    )
  }
}