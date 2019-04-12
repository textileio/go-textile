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
export default class Groups extends ConnectedComponent<RouteComponentProps, Stores> {
  state = {
    isLoading: false,
    isAdding: false,
    isRemoving: false,
    currentGroup: ''
  }
  handleRefreshClick = () => {
    this.stores.store.fetchGroups()
    this.setState({ isLoading: true })
    // Show spinner to indicate work is being done
    setTimeout(() => this.setState({ isLoading: false}), 3000)
  }
  handleAddGroup = (data: any) => {
  }
  handleShowRemove = (id: string) => {
    this.setState({ isRemoving: true, currentGroup: id })
  }
  handleDone = () => this.setState({ currentGroup: '', isRemoving: false })
  handleConfirm = () => {
    // this.stores.store.removeGroup(this.state.currentGroup).then(this.handleDone)
  }
  render() {
    const { groups } = this.stores.store
    console.log(groups.map((item: any) => { return {...item} }))
    return (
      <div style={{ height: '100vh' }}>
        <Segment basic style={{ height: 'calc(100vh-50px)'}}>
          <Header as='h3'>
            GROUPS
          <Header.Subheader>
              View and edit Groups
            </Header.Subheader>
          </Header>
          <Card.Group>
            {groups && groups.map((item: any) => this.renderItem(item))}
          </Card.Group>
        </Segment>
        <Button.Group fluid widths='2' style={{ position: 'absolute', bottom: 0 }}>
          <Button
            // disabled={groups.length < 1}
            style={{ borderRadius: 0 }}
            loading={this.state.isLoading}
            onClick={this.handleRefreshClick}
            content='Refresh' icon='refresh'
            positive type='button' />
          {/* <AddGroupModal
            open={this.state.isAdding}
            onClose={() => { this.setState({isAdding: false}) }}
            handleCafeAdd={this.handleAddGroup} trigger={
            <Button
              onClick={() => { this.setState({ isAdding: true }) }}
              style={{ borderRadius: 0 }}
              content='Add' icon='plus' type='button'/>
          }/> */}
        </Button.Group>
        <BackArrow onClick={() => { this.props.navigate && this.props.navigate('..') }} />
      </div>
    )
  }
  renderItem(item: any) {
    console.log(item)
    return (
      <Card key={item.id}>
        <Card.Content>
          <Card.Header>{item.name}</Card.Header>
          <Card.Meta>
            {item.type.toLowerCase().replace('_', ' ')},
            {item.sharing.toLowerCase().replace('_', ' ')}
          </Card.Meta>
          <Card.Description>
            {item.peer_count} Peers sharing {item.block_count} items
          </Card.Description>
        </Card.Content>
        <Card.Content extra>
          Updated <Moment fromNow>{item.head_block.date}</Moment>
        </Card.Content>
        <Icon
          style={{ position: 'absolute', right: '5px', top: '5px' }}
          link name='close'
          onClick={() => { this.handleShowRemove(item.id) }}
        />
      </Card>
    )
  }
}
