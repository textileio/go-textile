import React, { createRef, SyntheticEvent } from 'react'
import { observer } from 'mobx-react'
import { Segment, Label, Icon, Image, Input, Form, Header, Button, Card } from 'semantic-ui-react'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { RouteComponentProps } from '@reach/router'
import BackArrow from '../Components/BackArrow'
import AddCafeModal from '../Components/AddCafeModal'
import { Stores } from '../Stores'
import Moment from 'react-moment'
const { shell } = window.require('electron')

@connect('store') @observer
export default class Cafes extends ConnectedComponent<RouteComponentProps, Stores> {
  state = {
    isLoading: false,
    isAdding: false
  }
  handleMessagesClick = () => {
    this.stores.store.fetchMessages()
    this.setState({ isLoading: true })
    // Show spinner to indicate work is being done
    setTimeout(() => this.setState({ isLoading: false}), 3000)
  }
  handleAddCafe = (data: any) => {
    const { url, token } = data
    this.stores.store.addCafe(url, token)
    this.setState({isAdding: false})
  }
  handleRemoveCafe = (id: string) => {
    this.stores.store.removeCafe(id)
  }
  render() {
    const { cafes } = this.stores.store
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic style={{ height: 'calc(100vh-50px)'}}>
          <Header as='h3'>
            CAFES
          <Header.Subheader>
              Add and remove Cafes
            </Header.Subheader>
          </Header>
          <Card.Group>
            {cafes && cafes.map((item: any) => this.renderItem(item))}
          </Card.Group>
        </Segment>
        <Button.Group fluid widths='2' style={{ position: 'absolute', bottom: 0 }}>
          <Button
            disabled={cafes.length < 1}
            style={{ borderRadius: 0 }}
            loading={this.state.isLoading}
            onClick={this.handleMessagesClick}
            content='Messages' icon='refresh'
            positive type='button' />
          <AddCafeModal
            open={this.state.isAdding}
            onClose={() => { this.setState({isAdding: false}) }}
            handleCafeAdd={this.handleAddCafe} trigger={
            <Button
              onClick={() => { this.setState({ isAdding: true }) }}
              style={{ borderRadius: 0 }}
              content='Add' icon='plus' type='button'/>
          }/>
        </Button.Group>
        <BackArrow onClick={() => { this.props.navigate && this.props.navigate('..') }} />
      </div>
    )
  }
  renderItem(item: any) {
    const { cafe } = item
    return (
      <Card key={item.id}>
        <Card.Content>
          <Card.Header>{cafe.address.slice(0, 8)}</Card.Header>
          <Card.Meta>API {cafe.api}, Node v{cafe.node}</Card.Meta>
          <Card.Description
            as='a'
            onClick={() => { shell.openExternal(cafe.url) }}
            >{cafe.url}</Card.Description>
        </Card.Content>
        <Card.Content extra>
          Token expires <Moment fromNow>{item.exp}</Moment>
        </Card.Content>
        <Icon
          style={{ position: 'absolute', right: '5px', top: '5px' }}
          link name='close'
          onClick={() => { this.handleRemoveCafe(item.id) }}
        />
      </Card>
    )
  }
}
