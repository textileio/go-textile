import { Clipboard, Shell, Remote } from 'electron'

declare global {
  interface Window {
    require: (module: 'electron') => {
      clipboard: Clipboard,
      shell: Shell,
      remote: Remote
    }
  }
}