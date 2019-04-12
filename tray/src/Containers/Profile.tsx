import React, { createRef, SyntheticEvent } from 'react'
import { observer } from 'mobx-react'
import { Segment, Label, Icon, Image, Input, Form, Header, Button } from 'semantic-ui-react'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { RouteComponentProps } from '@reach/router'
import BackArrow from '../Components/BackArrow'
import { Stores } from '../Stores'
import Moment from 'react-moment'
const { clipboard } = window.require('electron')

@connect('store') @observer
export default class Profile extends ConnectedComponent<RouteComponentProps, Stores> {
  state = {
    isLoading: false
  }
  handleUsername = (e: SyntheticEvent) => {
    e.preventDefault()
    const current = this.inputRef.current
    if (current) {
      this.stores.store.setProfile(current.value, undefined)
    }
  }
  handleAvatar = (e: SyntheticEvent) => {
    e.preventDefault()
    const files = (e.target as HTMLInputElement).files
    if (files && files.length > 0) {
      const form = new FormData()
      form.append('file', files[0], files[0].name)
      this.stores.store.setProfile(undefined, form)
    }
  }
  private inputRef = createRef<HTMLInputElement>()
  private fileUploader = createRef<HTMLInputElement>()
  onAddressClick = () => {
    const { profile } = this.stores.store
    if (profile) {
      clipboard.write({ text: profile.address })
    }
  }
  onCafesClick = () => {
    this.stores.store.fetchCafes().then(() => {
      this.props.navigate && this.props.navigate('./cafes')
    })
  }
  onGroupsClick = () => {
    this.stores.store.fetchGroups().then(() => {
      this.props.navigate && this.props.navigate('./groups')
    })
  }
  handleAccountSync = () => {
    this.stores.store.syncAccount()
    this.setState({ isLoading: true })
    // Show spinner to indicate work is being done
    setTimeout(() => this.setState({ isLoading: false }), 3000)
  }
  handleLogout = () => {
    this.props.navigate && this.props.navigate('/landing')
  }
  render() {
    const { profile } = this.stores.store
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic style={{ height: '100vh' }}>
          <Header as='h3' onClick={this.onAddressClick}>
            ACCOUNT
            <Header.Subheader>
              Updated {profile ? <Moment fromNow>{profile.date}</Moment> : 'never'}
            </Header.Subheader>
          </Header>
          <input
            type="file"
            id="file"
            ref={this.fileUploader}
            style={{ display: "none" }}
            onChange={this.handleAvatar}
          />
          
            <Segment basic style={{ padding: 0 }}>
              {profile &&
                <Image
                  style={{ objectFit: 'cover', width: '150px', height: '150px' }}
                  centered
                  circular
                  src={profile.avatar} size='small'
                />
              }
            <Label as='a' style={{ padding: '1em 0 0 1em', border: 'none', left: '20%' }}
                basic attached='top left' size='large'
                onClick={() => {
                  if (this.fileUploader.current) {
                    this.fileUploader.current.click()
                  }
                }}
              >
                <Icon.Group style={{ margin: 0 }} size='large'>
                  <Icon name='image outline' />
                  <Icon corner='top right' name='pencil' />
                </Icon.Group>
              </Label>
            </Segment>
          <Header as='h4' style={{ margin: '1em 0 0.2em 0'}}>USERNAME</Header>
          <Form onSubmit={this.handleUsername}>
            <Form.Field>
              <Input
                iconPosition='left'
                labelPosition='right'
                placeholder='username'
                defaultValue={profile ? profile.name : ''}
              >
                <Icon name='pencil' />
                <input ref={this.inputRef} />
                <Label icon='save' onClick={this.handleUsername}/>
              </Input>
            </Form.Field>
          </Form>
          <Header as='h4' style={{ margin: '1em 0 0.2em 0' }}>INFO</Header>
          <Button.Group basic fluid compact>
            <Button content='Cafes' icon='coffee' type='button' onClick={this.onCafesClick} />
            <Button content='Groups' icon='users' type='button' onClick={this.onGroupsClick}/>
          </Button.Group>
        </Segment>
        <Button.Group fluid widths='2' style={{ position: 'absolute', bottom: 0 }}>
          <Button
            style={{ borderRadius: 0 }}
            loading={this.state.isLoading}
            content='Sync' icon='refresh' positive type='button'
            onClick={this.handleAccountSync}/>
          <Button
            disabled
            style={{ borderRadius: 0 }}
            content='Log-out' icon='log out' type='button'
            onClick={this.handleLogout} />
        </Button.Group>
        <BackArrow onClick={() => { this.props.navigate && this.props.navigate('..') }} />
      </div>
    )
  }
}
