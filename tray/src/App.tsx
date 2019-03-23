import React, { Component } from 'react'
import { observer } from 'mobx-react'
import 'react-semantic-toasts/styles/react-semantic-alert.css'
import Main from './Main'
import Splash from './Splash'
import { Dimmer, Loader } from 'semantic-ui-react'
import { SemanticToastContainer } from 'react-semantic-toasts'
import { Stores } from './Store'
import { ConnectedComponent, connect } from './ConnectedComponent'

interface AppProps { }

@connect('store') @observer
class App extends ConnectedComponent<AppProps, Stores> {
  componentDidMount() {
    const { store } = this.stores
     setTimeout(() => { store.checkStatus() }, 3000)
  }
  render() {
    const { store } = this.stores
    const view = ((screen: string) => {
      switch (screen) {
        case 'loading':
          return (
            <Splash />
          )
        case 'online':
          return (
            <Main />
          )
        default:
          return (
            <Dimmer inverted active>
              <Loader size='massive' />
            </Dimmer>
          )
      }
    })(store.status)
    return (
      <div className='App'>
        {view}
        <SemanticToastContainer />
      </div>
    )
  }
}

export default App
