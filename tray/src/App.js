import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import 'react-semantic-toasts/styles/react-semantic-alert.css'
import Main from './Main'
import { Dimmer, Loader } from 'semantic-ui-react'
import { SemanticToastContainer } from 'react-semantic-toasts'

@inject('store') @observer
class App extends Component {
  componentDidMount() {
    const { store } = this.props
    store.checkStatus()
  }
  render() {
    const { store } = this.props
    const view = (screen => {
      switch (screen) {
        case 'online':
          return (
            <div style={{ width: '80%', maxWidth: '500px', margin: '1em auto' }}>
              <Main />
            </div>
          )
        default:
          return (
            <Dimmer active={store.status !== 'online'}>
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
