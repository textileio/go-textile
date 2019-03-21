import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Image } from 'semantic-ui-react'
// import Moment from 'react-moment'
import * as Logo from './LaunchLogo@3x.png'

@inject('store') @observer
class Main extends Component {
  render() {
    const { store } = this.props
    return (
      <div>
        <Image centered size='medium' src={Logo} />
        <p>
          {store.profile.address}
        </p>
      </div>
    )
  }
}

export default Main
