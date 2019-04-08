import { observe, action, observable } from 'mobx'
import { Textile, Wallet } from '@textile/js-http-client'
import { createMemorySource, createHistory } from "@reach/router"

const textile = new Textile({
  url: 'http://127.0.0.1',
  port: 40602
})

export interface Message {
  name: string
  payload?: any
}

const source = createMemorySource("/")
const history = createHistory(source)

export interface Store {}

export class AppStore implements Store {
  constructor() {
    observe(this, 'status', change => {
      history.navigate(`/${change.newValue}`, { replace: false })
      if (change.newValue === 'online') {
        textile.profile.get().then((profile) => {
          if (!profile.name) {
            profile.name = profile.address.slice(-8)
          }
          if (profile.avatar) {
            profile.avatar = `${this.gateway}/ipfs/${profile.avatar}/0/small/d`
          } else {
            // TODO: Find a more permanent solution
            profile.avatar = 'https://react.semantic-ui.com/images/wireframe/square-image.png'
          }
          this.profile = profile
          textile.account.seed().then((value) => {
            this.accountSeed = value
          })
        }).catch((err: Error) => {
          console.log(err)
        })
      }
    })
    if ('astilectron' in window) {
      astilectron.onMessage((message: Message) => {
        switch (message.name) {
          default:
            console.log(message)
        }
      })
    }
  }
  createMnemonic(): string {
    return Wallet.fromWordCount(12).recoveryPhrase
  }
  async sendMessage(message: Message): Promise<Message> {
    message.payload = message.payload !== undefined ? message.payload : {}
    return new Promise(resolve => {
      astilectron.sendMessage(message, (response: Message) => {
        this.latestMessage = response
        resolve(response)
      })
    })
  }
  @action async initAndStartTextile(mnemonic: string, password: string) {
    if ('astilectron' in window) {
      try {
        const response = await this.sendMessage({
          name: 'init',
          payload: { mnemonic, password }
        })
        if (response) {
          this.status = 'online'
        } else {
          this.status = 'error'
        }
      } catch(err) {
        console.log(err)
        this.status = 'error'
      }
    } else {
      // Do nothing (we're probably in dev mode?)
      this.status = 'online'
    }
  }
  @observable accountSeed = ''
  @observable latestMessage: any = {}
  @observable history = history
  @observable gateway = 'http://127.0.0.1:5052'
  @observable status = 'starting'
  // TODO: Get proper types from js-http-client when Typescript lands
  @observable profile: any = {}
  @action checkStatus() {
    textile.utils.online().then((online: boolean) => {
      if (online) {
        this.status = 'online'
      } else {
        this.status = 'offline'
      }
    }).catch((err) => {
      console.log(err)
    })
  }
}

export interface Stores {
  store: AppStore
}
