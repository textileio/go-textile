import { observe, action, observable } from 'mobx'
import { Textile } from '@textileio/js-http-client'
import { toast } from 'react-semantic-toasts'

const textile = new Textile({
  url: 'http://127.0.0.1',
  port: 40602
})

class Store {
  constructor() {
    observe(this, 'status', change => {
      if (change.newValue === 'online') {
        textile.profile.get().then(profile => {
          if (!profile.username) {
            profile.username = profile.address.slice(-8)
          }
          this.profile = profile
        }).catch(err => {
          console.log(err)
        })
      }
    })
  }
  @observable gateway = 'http://127.0.0.1:5052'
  @observable status = 'offline'
  @observable profile = {}
  @action checkStatus() {
    textile.profile.get().then(online => {
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

export default Store
