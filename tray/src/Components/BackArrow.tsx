import React, { Component } from 'react'
import { Icon, IconProps } from 'semantic-ui-react'

export default class BackArrow extends Component<IconProps> {
  render() {
    const { name, ...props } = this.props
    return (
      <Icon {...props}
        style={{ position: 'absolute', right: '5px', top: '5px', zIndex: '1001' }}
        link name={name ? name : 'arrow left'}
      />
    )
  }
}
