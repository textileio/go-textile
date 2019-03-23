declare module 'react-semantic-toasts' {
  import { ReactNode } from 'react'
  import { SemanticICONS } from 'semantic-ui-react'
  type ContainerPosition = 'top-right' | 'top-center' | 'top-left' | 'bottom-right' | 'bottom-center' | 'bottom-left' 
  const SemanticToastContainer: (
    props: {
      position?: ContainerPosition
      animation?: string
    }
  ) => JSX.Element
  const toast: (
    options: {
      title: string
      description: string
      type?: 'info' | 'success' | 'warning' | 'error'
      icon?: SemanticICONS
      time?: number
    },
    onClose?: () => void,
    onClick?: () => void,
  ) => void

  export { SemanticToastContainer, toast };
}
