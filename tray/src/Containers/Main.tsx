import React from 'react'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import Notifications from '../Components/Notifications'
import { Image, Segment, Menu, Dropdown, Icon } from 'semantic-ui-react'
import { RouteComponentProps } from '@reach/router'
import { observer } from "mobx-react"
import { Stores } from '../Stores'
const { shell } = window.require('electron')

@connect('store') @observer
export default class Summary extends ConnectedComponent<RouteComponentProps, Stores> {
  onAPIClick = () => { shell.openExternal('http://127.0.0.1:40602/docs/index.html') }
  onQuitClick = () => {
    this.stores.store.sendMessage({ name: 'quit' })
  }
  onAccountClick = () => {
    this.props.navigate && this.props.navigate('./profile')
  }
  onNotificationsClick = () => {
    this.stores.store.fetchNotifications()
  }
  render() {
    const { store } = this.stores
    return (
      <div>
        <Menu attached='top' borderless>
          <Menu.Item header as='h3' style={{ padding: '10px' }}>
            {store.profile && <Image avatar src={store.profile.avatar}/>}
            {store.profile && store.profile.name }
          </Menu.Item>
          <Menu.Menu position='right'>
            <Dropdown item icon={{size: 'large', name: 'setting'}}>
              <Dropdown.Menu>
                <Dropdown.Item icon='user' text='Account' onClick={this.onAccountClick} />
                <Dropdown.Item icon='wrench' text='Settings' disabled/>
                <Dropdown.Item icon='refresh' text='Notifications' onClick={this.onNotificationsClick} />
                <Dropdown.Divider />
                <Dropdown.Item icon='external' text='API Docs' onClick={this.onAPIClick}/>
                <Dropdown.Divider />
                <Dropdown.Item icon='close' text='Quit' onClick={this.onQuitClick}/>
              </Dropdown.Menu>
            </Dropdown>
          </Menu.Menu>
        </Menu>
        <Segment basic style={{ height: '88vh', overflowY: 'auto', marginTop: 0 }} >
          <Notifications />
        </Segment>
      </div>
    )
  }
}
