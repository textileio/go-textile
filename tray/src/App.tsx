import React from 'react'
import { observer } from 'mobx-react'
import 'react-semantic-toasts/styles/react-semantic-alert.css'
import { Router, LocationProvider } from '@reach/router'
import Login from './Components/Login'
import Profile from './Components/Profile'
import Create from './Components/Create'
import Summary from './Components/Summary'
import Basic from './Components/Basic'
import Splash from './Components/Splash'
import Start from './Components/Start'
import { Stores } from './Stores'
import { ConnectedComponent, connect } from './Components/ConnectedComponent'

interface AppProps { }

@connect('store') @observer
class App extends ConnectedComponent<AppProps, Stores> {
  componentDidMount() {
    const { store } = this.stores
  //   setTimeout(() => { store.status = 'onboard' }, 3000)
    store.status = 'online'
  }
  render() {
    const { store } = this.stores
    return (
      <LocationProvider history={store.history}>
        <Router>
          <Splash default />
          <Basic path='/onboard'>
            <Start path='/'/>
            <Login path='/login' />
            <Create path='/create' />
          </Basic>
          <Splash path='/loading' />
          <Basic path='/online'>
            <Summary path='/' />
            <Profile path='/profile' />
          </Basic>
        </Router >
      </LocationProvider>
    )
  }
}

export default App
