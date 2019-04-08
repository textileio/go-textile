import React from 'react'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { Image, Segment, Menu, Icon, Dropdown } from 'semantic-ui-react'
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
  render() {
    const { store } = this.stores
    return (
      <div>
        <Menu attached='top' borderless>
          <Menu.Item header>
            <Image avatar src={store.profile.avatar}/>
            {store.profile.name}
          </Menu.Item>
          <Menu.Menu position='right'>
            <Dropdown item icon='setting'>
              <Dropdown.Menu>
                <Dropdown.Item icon='user' text='Account' onClick={this.onAccountClick} />
                <Dropdown.Item icon='wrench' text='Settings' />
                <Dropdown.Divider />
                {/* <Dropdown.Header icon='code' content='Develop' /> */}
                <Dropdown.Item icon='external' text='API Docs' onClick={this.onAPIClick}/>
                <Dropdown.Divider />
                <Dropdown.Item icon='close' text='Quit' onClick={this.onQuitClick}/>
              </Dropdown.Menu>
            </Dropdown>
          </Menu.Menu>
        </Menu>
        <Segment basic attached='bottom'>
          Content here
        </Segment>
      </div>
    )
  }
}
