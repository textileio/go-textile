import { observe, action, observable } from 'mobx'
import { Textile } from '@textileio/js-http-client'
import { toast } from 'react-semantic-toasts'

const textile = new Textile({
  url: 'http://127.0.0.1',
  port: 40602
})

export interface Store {}

export class AppStore implements Store {
  constructor() {
    observe(this, 'status', change => {
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
  }
  @observable gateway: string = 'http://127.0.0.1:5052'
  @observable status: string = 'loading'
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

export const stores: Stores = {
  store: new AppStore()
}
