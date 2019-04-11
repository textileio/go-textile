import React, { Component, SyntheticEvent } from 'react'
import { Modal, Form, Button, ModalProps, Input, Message } from 'semantic-ui-react'

interface State {
  url: string
  token: string
}

interface Props {
  handleCafeAdd: (state: State) => void
}

const validURL = (str: string) => {
  var pattern = new RegExp('^(https?:\\/\\/)?' + // protocol
    '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|' + // domain name
    '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
    '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*' + // port and path
    '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
    '(\\#[-a-z\\d_]*)?$', 'i') // fragment locator
  return pattern.test(str)
}

export default class AddCafeModal extends Component<ModalProps & Props, State> {
  state = {
    url: '',
    token: ''
  }
  handleSubmit = () => {
    this.props.handleCafeAdd(this.state)
    this.setState({ url: '', token: '' })
  }
  render() {
    const { handleCafeAdd, ...props } = this.props
      return (
        <Modal {...props}
          // closeIcon={{ name: 'close', color: 'black' }}
          dimmer='inverted'
          size='tiny'
        >
        {/* <Modal.Header>Add a Cafe</Modal.Header> */}
          <Modal.Content>
            <Modal.Description>
              <Form onSubmit={this.handleSubmit}>
              <Form.Field
                title='Enter a valid URL (full url or IP address)'
                required error={!validURL(this.state.url)}>
              <label>URL</label>
              <Input
                autoFocus
                name='url'
                value={this.state.url}
                placeholder={'Cafe\'s public address...'}
                onChange={(e: SyntheticEvent, data: any) => this.setState({ url: data.value })}
              />
              </Form.Field>
                <Form.Field
                  title='Enter a valid access token (60 chars)'
                  required error={this.state.token.length !== 60}>
                <label>ACCESS TOKEN</label>
                <Input
                  name='token'
                  value={this.state.token}
                  placeholder='Required token...'
                  onChange={(e: SyntheticEvent, data: any) => this.setState({ token: data.value })}
                />
              </Form.Field>
                <Button 
                  disabled={!validURL(this.state.url) || this.state.token.length !== 60}
                  content='Add' icon='add circle' type='submit' />
            </Form>
          </Modal.Description>
        </Modal.Content>
      </Modal>
    )
  }
}
