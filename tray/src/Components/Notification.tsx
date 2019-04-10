import React, { Component } from 'react'
import { Feed, Icon } from 'semantic-ui-react'
import { observer } from 'mobx-react'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { Stores } from '../Stores'
import Moment from 'react-moment'

interface Props { 
  item: any,
}

@connect('store') @observer
export default class Notification extends ConnectedComponent<Props, Stores> {
  state = {
    showClose: false
  }
  toggleCloseVisible = () => {
    this.setState({ showClose: !this.state.showClose })
  }
  handleClose = () => {
    this.stores.store.readNotification(this.props.item.id)
  }
  render() {
    const { item } = this.props
    const isMessage = item.type === 'MESSAGE_ADDED'
    return (
      <Feed.Event
        onMouseEnter={() => this.toggleCloseVisible()}
        onMouseLeave={() => this.toggleCloseVisible()}
        style={{ marginBottom: '0.5em' }}>
        {/* <Feed.Label>
            <img src={item.user.avatar} />
          </Feed.Label> */}
        <Feed.Content>
          <Feed.Date style={{ height: '1.5em' }}>
            {this.state.showClose &&
            <div>
              <Moment fromNow>{item.date}</Moment>
              {!item.read &&
                <Icon title='mark as read' link hidden={true} name='check' onClick={this.handleClose}
                  style={{ float: 'right' }}
                />
              }
            </div>
            }
          </Feed.Date>
          <Feed.Summary style={{ fontWeight: 'normal', color: item.read ? 'lightgrey' : 'black' }}>
            {item.user.name} {isMessage ? 'posted' : item.body} in <span style={{ fontWeight: 'bold' }}>{item.subject_desc}</span>
          </Feed.Summary>
          {isMessage && <Feed.Meta style={{ marginTop: 0, color: item.read ? 'lightgrey' : 'black' }}>{item.body}</Feed.Meta>}
        </Feed.Content>
      </Feed.Event>
    )
  }
}
