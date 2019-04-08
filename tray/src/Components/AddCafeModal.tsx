import React, { Component } from 'react'
import { Modal, Form, Button, Image, TextArea } from 'semantic-ui-react'

export default class AddCafeModal extends Component {
  state = {
    url: '',
    token: ''
  }
  handleSubmit = () => {
    // this.props.onSubmit(this.state)
  }
  handleChange = (e, { name, value }) => this.setState({ [name]: value })
  render() {
    const { url, token } = this.state
  
    return (
      <Modal
        open={open}
        closeIcon={{ name: 'close', color: 'black' }}
        dimmer='inverted'
        onClose={onClose}
        size='tiny'
      >
        <Modal.Header>Add a file</Modal.Header>
        <Modal.Content>
          <Modal.Description>
            <Form onSubmit={this.handleSubmit}>
              <Form.Field>
                <TextArea
                  autoFocus
                  name='caption'
                  label='Caption'
                  value={caption}
                  placeholder='Add a caption...'
                  // onChange={this.handleChange}
                />
                <Button>Add</Button>
              </Form.Field>
            </Form>
          </Modal.Description>
        </Modal.Content>
      </Modal>
    )
  }
}
