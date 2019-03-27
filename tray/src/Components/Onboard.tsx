import React, { ReactElement } from 'react'
import { RouteComponentProps } from '@reach/router'

const Onboard = (props: { children?: ReactElement[] } & RouteComponentProps) => <div>{props.children}</div>
export default Onboard
