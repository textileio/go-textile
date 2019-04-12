import React, { SyntheticEvent } from 'react'
import { Button, Header, Image, Segment, Form, Icon } from 'semantic-ui-react'
import { RouteComponentProps } from '@reach/router'
import { DropdownItemProps } from 'semantic-ui-react'
import { ConnectedComponent, connect } from '../Components/ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'
import BackArrow from '../Components/BackArrow'
import NotificationsImage from '../Assets/notifications@3x.png'

interface State {
  address: string
  password: string
  passType: string
}

@connect('store') @observer
export default class Landing extends ConnectedComponent<RouteComponentProps, Stores, State> {
  state = {
    address: '',
    password: '',
    passType: 'password'
  }
  componentDidMount() {
    const { addresses } = this.stores.store
    this.setState({ address: addresses.length > 0 ? addresses[0] : '' })
  }
  handleCreate = () => this.props.navigate && this.props.navigate('/onboard')
  handleLogin = () => {
    this.stores.store.screen = 'loading'
    this.stores.store.initAndStartTextile(undefined, this.state.address, this.state.password)
  }
  togglePassType = () => this.setState({
    passType: this.state.passType === 'password' ? 'input' : 'password'
  })
  render() {
    const { addresses } = this.stores.store
    const options = addresses.map((item: string) => {
      const text = item.slice(0, 8) + '...' + item.slice(-8)
      const data: DropdownItemProps = { text, value: item, key: item }
      return data
    })
    return (
      <div>
        <Form onSubmit={this.handleLogin} style={{ height: '100vh' }}>
          <Segment basic>
            <Image centered size='small' src={NotificationsImage} />
            <Header as='h3'>
              Welcome back!
              <Header.Subheader>
                Log-in using an existing account, or create a new one.
              </Header.Subheader>
            </Header>
            <Form.Field title='Choose existing account'>
              <label>ACCOUNTS</label>
              <Form.Select
                autoFocus
                name='address'
                placeholder='Account Address'
                options={options}
                defaultValue={addresses.length > 0 ? addresses[0] : ''}
                onChange={(e: SyntheticEvent, data: any) => this.setState({ address: data.value })}
              />
            </Form.Field>
            <Form.Field title={'Enter your password (can be left blank)'}>
              <label>PASSWORD</label>
              <Form.Input
                name='password'
                type={this.state.passType}
                value={this.state.password}
                placeholder='Password...'
                icon={<Icon
                  name={this.state.passType === 'password' ? 'eye' : 'eye slash'}
                  link
                  onClick={this.togglePassType}
                />}
                onChange={(e: SyntheticEvent, data: any) => this.setState({ password: data.value })}
              />
            </Form.Field>
          </Segment>
          <Button.Group fluid widths='2' style={{ position: 'absolute', bottom: 0 }}>
            <Button style={{ borderRadius: 0 }} positive content='Sign-in' icon='sign-in' type='submit' />
            <Button style={{ borderRadius: 0 }} content='New' icon='key' type='button' onClick={this.handleCreate} />
          </Button.Group>
        </Form>
        <BackArrow name='close' onClick={() => { this.stores.store.sendMessage({ name: 'quit' }) }} />
      </div>
    )
  }
}