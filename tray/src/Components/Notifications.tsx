import React from 'react'
import { observer } from 'mobx-react'
import { Feed } from 'semantic-ui-react'
import { ConnectedComponent, connect } from './ConnectedComponent'
import Notification from './Notification'
import { Stores } from '../Stores'

interface Props {}

@connect('store') @observer
export default class Notifications extends ConnectedComponent<Props, Stores> {
  state = {
    isLoading: false
  }
  render() {
    const { notifications } = this.stores.store
    return (
      <Feed size='small'>
        {notifications && notifications.map((item: any) => <Notification key={item.id} item={item} />)}
      </Feed>
    )
  }
}
