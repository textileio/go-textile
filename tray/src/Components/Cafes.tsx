import React, { createRef } from 'react'
import { observer } from 'mobx-react'
import { Segment, Label, Icon, Image, Input, Form, Header, Button, Card } from 'semantic-ui-react'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { RouteComponentProps } from '@reach/router'
import BackArrow from './BackArrow'
import { Stores } from '../Stores'
import Moment from 'react-moment'
const { clipboard } = window.require('electron')

@connect('store') @observer
export default class Cafes extends ConnectedComponent<RouteComponentProps, Stores> {
  handleAddCafe = (e: any) => {
    e.preventDefault()
    // this.stores.store.setProfile(this.inputRef.value, null)
  }
  handleRemoveCafe = (e: any) => {
    e.preventDefault()
  //   this.stores.store.setProfile(null, form)
  }
  private inputRef = createRef<HTMLInputElement>()
  onAddressClick = () => { clipboard.write({ text: this.stores.store.profile.address }) }
  onSeedClick = () => { clipboard.write({ text: this.stores.store.accountSeed }) }
  render() {
    const { profile } = this.stores.store
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic style={{ height: '100vh', overflowY: 'scroll' }}>
          <Header as='h3'>
            ACCOUNT
          <Header.Subheader>
              Updated <Moment fromNow>{profile.updated}</Moment>
            </Header.Subheader>
          </Header>
          <input
            type="file"
            id="file"
            ref="fileUploader"
            style={{ display: "none" }}
          // onChange={this.handleAvatar}
          />
          <Form onSubmit={this.handleUsername}>
            <Form.Field style={{ textAlign: 'center' }}>
              <Segment circular basic style={{ display: 'inline-block', padding: 0 }}>
                <Image centered circular src={profile.avatar} size='small'
                // onClick={() => { this.refs.fileUploader.click() }}
                />
                <Label circular as='a' basic attached='top left' size='large'>
                  <Icon.Group style={{ margin: 0 }} size='large'>
                    <Icon name='image outline' />
                    <Icon corner='top right' name='add' />
                  </Icon.Group>
                </Label>
              </Segment>
            </Form.Field>
            <Form.Field>
              <label>USERNAME</label>
              <Input iconPosition='left' defaultValue={profile.name}>
                <Icon name='pencil' />
                <input ref={this.inputRef} />
              </Input>
            </Form.Field>
            <Form.Field>
              <label>INFO</label>
              <Button.Group basic fluid compact>
                <Button content='Cafes' icon='coffee' type='button' />
                <Button content='Address' icon='copy outline' type='button' onClick={this.onAddressClick}/>
                <Button content='Seed' icon='copy outline' type='button' onClick={this.onSeedClick}/>
              </Button.Group>
            </Form.Field>
          </Form>
          <Card>
            <Card.Content>
              <Card.Header>Cafe name</Card.Header>
              <Card.Meta>API v0, Node v0.1.11, copy peer</Card.Meta>
              <Card.Description>https://cafe.textile.io/</Card.Description>
            </Card.Content>
            <Card.Content extra>
              Updated in the past
            </Card.Content>
          </Card>
        </Segment>
        <Button.Group fluid style={{ position: 'absolute', bottom: 0 }}>
          <Button style={{ borderRadius: 0 }} content='Sync' icon='refresh' positive type='button' />
          <Button style={{ borderRadius: 0 }} content='Log-out' icon='log out' type='button' />
        </Button.Group>
        <BackArrow onClick={() => { this.props.navigate && this.props.navigate('..') }} />
        <AddCafeModal
          onClose={this.handleClose}
          onSubmit={this.handleSubmit}
          open={this.state.modalOpen}
          preview={this.state.url}
        />
      </div>
    )
  }
}
