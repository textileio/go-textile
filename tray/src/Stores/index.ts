import { observe, action, observable, runInAction } from 'mobx'
import { Textile, Wallet } from '@textile/js-http-client'
import { createMemorySource, createHistory } from "@reach/router"
import moment, { utc } from 'moment'

const textile = new Textile({
  url: 'http://127.0.0.1',
  port: 40600
})

const AVATAR_KEY = 'profile'
const DEFAULT_AVATAR = 'https://react.semantic-ui.com/images/wireframe/square-image.png'

export interface Message {
  name: string
  payload?: any
}

const source = createMemorySource("/")
const history = createHistory(source)

export interface Store {}

export interface ProfileInfo {
  name: string
  avatar: string
  date: string
  address: string
}

export class AppStore implements Store {
  constructor() {
    observe(this, 'screen', change => {
      history.navigate(`/${change.newValue}`, { replace: false })
      switch(change.newValue) {
        case 'online':
          this.fetchProfile()
          this.fetchCafes()
          this.fetchNotifications()
          break
        default:
          break
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
  // Methods
  createMnemonic(): string {
    return Wallet.fromWordCount(12).recoveryPhrase
  }
  async sendMessage(message: Message): Promise<Message> {
    message.payload = message.payload !== undefined ? message.payload : {}
    return new Promise(resolve => {
      // TODO: Catch if astilectron doesn't exist
      astilectron.sendMessage(message, (response: Message) => {
        resolve(response)
      })
    })
  }
  // Observables
  @observable history = history
  @observable gateway = 'http://127.0.0.1:5050'
  @observable screen = 'starting'
  @observable cafes: any[] = []
  @observable notifications: any[] = []
  @observable profile?: ProfileInfo = undefined
  // Actions
  @action async initAndStartTextile(mnemonic: string, password: string) {
    let screen = 'starting'
    if ('astilectron' in window) {
      try {
        const response = await this.sendMessage({
          name: 'init',
          payload: { mnemonic, password }
        })
        if (response) {
          screen = 'online'
        } else {
          screen = 'error'
        }
      } catch(err) {
        screen = 'error'
      }
    } else {
      // Do nothing (we're probably in dev mode?)
      screen = 'online'
    }
    runInAction('initAndStartTextile', () => {
      this.screen = screen
    })
  }
  @action async fetchMessages() {
    try {
      const success = await textile.cafes.checkMessages()
    } catch(err) {
      console.log(err)
    }
  }
  @action async fetchProfile() {
    try {
      const contact = await textile.account.contact()
      let updated: string | undefined
      for (let peer of contact.peers) {
        if (moment(peer.updated).isAfter(moment(updated))) {
          updated = peer.updated
        }
      }
      runInAction('fetchProfile', () => {
        this.profile = {
          name: contact.name ? contact.name : contact.address.slice(-8),
          address: contact.address,
          avatar: contact.avatar ? `${this.gateway}/ipfs/${contact.avatar}/0/small/d` : DEFAULT_AVATAR,
          date: updated || utc().format()
        }
      })
    } catch(err) {
      console.log(err)
    }
  }
  @action async fetchNotifications() {
    try {
      const notifications = await textile.notifications.list()
      const processed = notifications.items.map((item) => {
        if (item.user.avatar) {
          item.user.avatar = `${this.gateway}/ipfs/${item.user.avatar}/0/small/d`
        } else {
          // TODO: Find a more permanent solution
          item.user.avatar = DEFAULT_AVATAR
        }
        return item
      })
      runInAction('fetchNotifications', () => {
        this.notifications = processed
      })
    } catch(err) {
      console.log(err)
    }
  }
  @action async readNotification(id: string) {
    try {
      const success = await textile.notifications.read(id)
      if (success) {
        this.fetchNotifications()
      }
    } catch(err) {
      console.log(err)
    }
  }
  @action async fetchCafes() {
    try {
      const cafes = await textile.cafes.list()
      runInAction('fetchCafes', () => {
        this.cafes = cafes.items
      })
    } catch(err) {
      console.log(err)
    }
  }
  @action async syncAccount() {
    textile.account.sync()
  }
  @action async setProfile(userString?: string, avatarFile?: FormData) {
    if (userString) {
      await textile.profile.setUsername(userString)
      this.fetchProfile()
    }
    if (avatarFile) {
      let avatarThread
      const threads = await textile.threads.list()
      for (const thread of threads.items) {
        if (thread.key === AVATAR_KEY) {
          avatarThread = thread
          break
        }
      }
      if (!avatarThread) {
        const schemas = await textile.schemas.defaults()
        const avatarSchema = schemas.avatar
        const file = await textile.schemas.add(avatarSchema)
          avatarThread = await textile.threads.add(
            'Profile Images', file.hash, AVATAR_KEY, 'private', 'not_shared'
          )
      }
      const addedFile = await textile.files.addFile(avatarFile, 'avatar', avatarThread.id)
      await textile.profile.setAvatar(addedFile.files[0].links.large.hash)
      this.fetchProfile()
    }
  }
  @action async fetchStatus() {
    try {
      const online = await textile.utils.online()
      if (online) {
        this.screen = 'online'
      } else {
        this.screen = 'offline'
      }
    } catch(err) {
      console.log(err)
    }
  }
  @action async addCafe(url: string, token: string) {
    try {
      await textile.cafes.add(url, token)
      this.fetchCafes()
    } catch (err) {
      console.log(err)
    }
  }
  @action async removeCafe(id: string) {
    try {
      const success = await textile.cafes.remove(id)
      if (success) {
        this.fetchCafes()
      } else {
        console.log("error!")
      }
    } catch (err) {
      console.log(err)
    }
  }
}

export interface Stores {
  store: AppStore
}
