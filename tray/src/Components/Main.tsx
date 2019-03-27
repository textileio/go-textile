import React from 'react'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

interface MainProps{
  path: string
}

@connect('store') @observer
class Main extends ConnectedComponent<MainProps, Stores> {
  render() {
    const { store } = this.stores
    return (
      <div>
        <p>{store.profile.address}</p>
      </div>
    )
  }
}

export default Main
