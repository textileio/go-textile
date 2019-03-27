import React from 'react'
import { Dimmer, Loader } from 'semantic-ui-react'

interface RouterProps {
  default?: boolean
  path?: string
}

const Loading = (props: RouterProps) => {
  return (
    <Dimmer inverted active>
      <Loader size='massive' />
    </Dimmer>
  )
}

export default Loading
