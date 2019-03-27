import React from 'react'
import { observer } from 'mobx-react'
import 'react-semantic-toasts/styles/react-semantic-alert.css'
import { SemanticToastContainer } from 'react-semantic-toasts'
import { Router, LocationProvider } from '@reach/router'
import { Grid } from 'semantic-ui-react'
import Onboard from './Components/Onboard'
import Login from './Components/Login'
import Create from './Components/Create'
import Loading from './Components/Loading'
import Success from './Components/Success'
import Splash from './Components/Splash'
import Start from './Components/Start'
import { Stores } from './Stores'
import { ConnectedComponent, connect } from './Components/ConnectedComponent'

interface AppProps { }

@connect('store') @observer
class App extends ConnectedComponent<AppProps, Stores> {
  componentDidMount() {
    const { store } = this.stores
    setTimeout(() => { store.status = 'onboard' }, 3000)
  }
  render() {
    const { store } = this.stores
    return (
      <Grid columns={1} style={{ height: 'calc(100vh - 2em)', margin: '2em 0 0 0' }}>
        <Grid.Row stretched>
          <Grid.Column>
            <LocationProvider history={store.history}>
              <Router>
                <Splash default />
                <Onboard path='/onboard'>
                  <Start path='/'/>
                  <Login path='/login' />
                  <Create path='/create' />
                </Onboard>
                <Loading path='/loading' />
                <Success path='/online' />
              </Router >
              <SemanticToastContainer />
            </LocationProvider>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    )
  }
}

export default App
