import React from 'react'
import { Feed, Image } from 'semantic-ui-react'
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
        onClick={this.handleClose}
        style={{ cursor: 'pointer', opacity: item.read ? 0.5 : 1 }}>
        <Feed.Label>
          <Image style={{ width: '2em', height: '2em' }} src={item.user.avatar} />
        </Feed.Label>
        <Feed.Content>
          <Feed.Summary>
            {item.subject_desc}
          </Feed.Summary>
          <Feed.Extra text style={{ margin: 0 }}>
            {`${item.user.name} ${isMessage ? 'said:' : ''} ${item.body} `}
            <Feed.Date as='span'>
              <Moment fromNow ago>{item.date}</Moment>
            </Feed.Date>
          </Feed.Extra>
        </Feed.Content>
      </Feed.Event>
    )
  }
}
