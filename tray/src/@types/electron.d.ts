import { Clipboard, Shell } from 'electron'

declare global {
  interface Window {
    require: (module: 'electron') => {
      clipboard: Clipboard,
      shell: Shell
    }
  }
}