import React, { SyntheticEvent, FormEvent, ChangeEvent } from 'react'
import {
  Button, Header, Segment, Form, TextAreaProps, Message, Icon, Modal, Popup,InputOnChangeData,
  Progress, Input, PopupProps
} from 'semantic-ui-react'
import { Slide } from 'react-reveal'
import zxcvbn from 'zxcvbn'
import { RouteComponentProps } from '@reach/router'
import QrReader from 'react-qr-reader'
import { ConnectedComponent, connect } from './ConnectedComponent'
import { observer } from "mobx-react"
import { Stores } from '../Stores'

interface LoginState {
  mnemonic: string
  password: string
  modalOpen: boolean
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
export default class Login extends ConnectedComponent<RouteComponentProps, Stores, LoginState> {
  state = {
    mnemonic: '',
    password: '',
    modalOpen: false,
    score: undefined,
    passType: 'password'
  }
  handlePassChange = (event: ChangeEvent, data: InputOnChangeData) => {
    if (data.name === 'password') {
      this.setState({ password: data.value as string, score: zxcvbn(data.value).score })
    }
  }
  handleMnemonicChange = (event: FormEvent, data: TextAreaProps) => {
    if (data.name === 'mnemonic') {
      this.setState({ mnemonic: data.value as string })
    }
  }
  handleSubmit = (event: SyntheticEvent) => {
    this.stores.store.status = 'loading'
    this.stores.store.initAndStartTextile(this.state.mnemonic, this.state.password)
  }
  handleScan = (data: string | null) => {
    if (data !== null) {
      this.setState({ mnemonic: data })
      this.handleQrClose()
    }
  }
  handleError = () => console.log("error")
  togglePassType = () => this.setState({
    passType: this.state.passType === 'password' ? 'input' : 'password'
  })
  handleQrOpen = () => this.setState({ modalOpen: true })
  handleQrClose = () => this.setState({ modalOpen: false })
  render() {
    const { mnemonic, password, passType, score } = this.state
    
    const inValid = mnemonic.split(/\b[^\s]+\b/).length < 13
    return (
      <Slide right opposite>
        <Segment raised>
          <Form onSubmit={this.handleSubmit}>
            <Segment basic attached>
              <Header as='h3'>
                Enter an existing <BIP39Popup trigger={<span style={{ textDecoration: 'underline' }}>mnemonic passphrase</span>} />
              </Header>
              <Form.TextArea
                style={{ fontSize: '1.6em' }}
                name='mnemonic'
                value={mnemonic}
                onChange={this.handleMnemonicChange}
                />
              <Message
                warning
                visible={inValid && mnemonic !== ''}
                header='Must be >12 words long'
                content={'Your mnemonic must be at least 12 words long.'}
              />
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
              <Button content='Sign-in' icon='sign-in' type='submit' positive disabled={inValid} />
              <Modal
                trigger={<Button content='Scan' icon='qrcode' type='button' onClick={this.handleQrOpen} />}
                open={this.state.modalOpen}
                onClose={this.handleQrClose}
                size='small'
                basic
              >
                <QrReader
                  onError={this.handleError}
                  onScan={this.handleScan}
                  style={{ width: '80%', maxWidth: '500px', margin: '0 auto' }}
                />
              </Modal>
            </Button.Group>
          </Form>
        </Segment>
      </Slide>
    )
  }
}
