import React, { SyntheticEvent, ChangeEvent, createRef } from 'react'
import { Button, Header, Segment, Form, Icon, Popup, InputOnChangeData, Progress, Input, PopupProps, TextArea, Ref } from 'semantic-ui-react'
import { Slide } from 'react-reveal'
import zxcvbn from 'zxcvbn'
import { RouteComponentProps } from '@reach/router'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

interface CreateState {
  mnemonic: string
  password: string
  score?: number
  passType: string // password | input
}

const BIP39Popup = (props: PopupProps) => {
  return (
    < Popup hoverable trigger = { props.trigger } >
      <Popup.Content>
        We use a <a href="https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki">BIP39</a> mnemonic.
        It cannot be recovered, so don't forget it!
      </Popup.Content>
    </Popup >
  )
}

const PasswordPopup = (props: PopupProps) => {
  return (
    < Popup hoverable trigger={props.trigger} >
      <Popup.Content>
        This password will be used to <a href="https://github.com/mutecomm/go-sqlcipher">encrypt</a> your local Textile database.
        It cannot be recovered, so don't forget it!
      </Popup.Content>
    </Popup >
  )
}

@connect('store') @observer
export default class Create extends ConnectedComponent<RouteComponentProps, Stores, CreateState> {
  private textArea = createRef<HTMLTextAreaElement>()
  constructor(props: RouteComponentProps) {
    super(props)
    this.state = {
      mnemonic: this.stores.store.createMnemonic(),
      password: '',
      score: undefined,
      passType: 'password'
    }
  }
  copyToClipboard = (event: SyntheticEvent) => {
    const text = this.textArea.current
    if (text) {
      text.select()
      document.execCommand('copy')
      const target = (event.target as HTMLButtonElement)
      target.focus()
    }
  }
  handleRefresh = () => {
    this.setState({
      mnemonic: this.stores.store.createMnemonic()
    })
  }
  handlePassChange = (event: ChangeEvent, data: InputOnChangeData) => {
    if (data.name === 'password') {
      this.setState({ password: data.value as string, score: zxcvbn(data.value).score })
    }
  }
  handleSubmit = (event: SyntheticEvent) => {
    this.stores.store.status = 'loading'
    this.stores.store.initAndStartTextile(this.state.mnemonic, this.state.password)
  }
  handleError = () => console.log("error")
  togglePassType = () => this.setState({
    passType: this.state.passType === 'password' ? 'input' : 'password'
  })
  render() {
    const { mnemonic, password, passType, score } = this.state
    return (
      <Slide right opposite>
        <Segment raised>
          <Form onSubmit={this.handleSubmit}>
            <Segment basic attached>
              <Header as='h3'>
                Here's your secret <BIP39Popup trigger={<span style={{ textDecoration: 'underline' }}>mnemonic passphrase</span>} />
              </Header>
              <Form.Field style={{ margin: 0 }}>
                <Ref innerRef={this.textArea}>
                  <TextArea
                    icon='search'
                    readOnly
                    style={{ fontSize: '1.6em' }}
                    name='mnemonic'
                    value={mnemonic}
                  />
                </Ref>
              </Form.Field>
              <Button.Group floated='right' basic size='mini' attached='bottom' style={{ marginTop: 0 }}>
                <Button content='Copy' icon='copy' type='button' onClick={this.copyToClipboard}/>
                <Button content='Refresh' icon='refresh' type='button' onClick={this.handleRefresh}/>
                </Button.Group>
                <Form.Field>
                <label>Use an <PasswordPopup trigger={<span style={{ textDecoration: 'underline' }}>additional password</span>} /> for added security</label>
                  <Input
                    name='password'
                    type={passType}
                    placeholder='Password...'
                    value={password}
                    onChange={this.handlePassChange}
                    icon={<Icon
                      name={passType === 'password' ? 'eye' : 'eye slash'}
                      link
                      onClick={this.togglePassType}
                    />}
                    />
                  <Progress attached='bottom' indicating value={score || 0} total={4} />
                </Form.Field>
            </Segment>
            <Button.Group attached='bottom'>
              <Button content='Create' icon='user secret' type='submit' positive />
              <Button content='Cancel' icon='cancel' disabled type='button' />
            </Button.Group>
          </Form>
        </Segment>
      </Slide>
    )
  }
}
