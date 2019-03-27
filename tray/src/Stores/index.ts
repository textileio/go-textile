import { observe, action, observable } from 'mobx'
import { Textile, Wallet } from '@textileio/js-http-client'
import { toast } from 'react-semantic-toasts'
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
        textile.profile.get().then((profile: any) => {
          if (!profile.username) {
            profile.username = profile.address.slice(-8)
          }
          this.profile = profile
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
        console.log(response)
        resolve(response)
      })
    })
  }
  @action async initAndStartTextile(mnemonic: string, password: string) {
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
    }
  }
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
    }).catch(() => {
      toast({
        title: 'Offline?',
        description: 'Looks like your Textile peer is offline 😔',
        time: 0
      })
    })
  }
}

export interface Stores {
  store: AppStore
}
