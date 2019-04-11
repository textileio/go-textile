import React, { ReactElement } from 'react'
import { observer } from 'mobx-react'
import { RouteComponentProps } from '@reach/router'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { Stores } from '../Stores'

@connect('store') @observer
export default class Basic extends ConnectedComponent<{ children?: ReactElement[] } & RouteComponentProps, Stores> {
  render() {
    const { store } = this.stores
    return (
      <div>
        {this.props.children}
      </div>
    )
  }
}
