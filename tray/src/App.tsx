import React from 'react'
import { observer } from 'mobx-react'
import { Router, LocationProvider } from '@reach/router'
import Login from './Containers/Login'
import Profile from './Containers/Profile'
import Create from './Containers/Create'
import Landing from './Containers/Landing'
import Main from './Containers/Main'
import Cafes from './Containers/Cafes'
import Basic from './Components/Basic'
import Splash from './Components/Splash'
import Start from './Containers/Start'
import { Stores } from './Stores'
import { ConnectedComponent, connect } from './Components/ConnectedComponent'

interface AppProps { }

@connect('store') @observer
class App extends ConnectedComponent<AppProps, Stores> {
  render() {
    const { store } = this.stores
    return (
      <LocationProvider history={store.history}>
        <Router>
          <Splash default />
          <Landing path='/landing' />
          <Basic path='/onboard'>
            <Start path='/'/>
            <Login path='/login' />
            <Create path='/create' />
          </Basic>
          <Basic path='/online'>
            <Main path='/' />
            <Basic path='/profile'>
              <Profile path='/' />
              <Cafes path='/cafes' />
            </Basic>
          </Basic>
        </Router >
      </LocationProvider>
    )
  }
}

export default App
