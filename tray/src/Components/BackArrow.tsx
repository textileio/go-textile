import React, { Component } from 'react'
import { Label, Icon, LabelProps } from 'semantic-ui-react'

export default class BackArrow extends Component<LabelProps> {
  render() {
    const { name, ...props } = this.props
    return (
      <Label {...props}
        as='a' style={{ position: 'absolute', right: '0', top: '0', zIndex: '10' }}
        basic size='large'
      >
        <Icon style={{ margin: 0 }} name={name ? name : 'arrow left'} />
      </Label>
    )
  }
}
