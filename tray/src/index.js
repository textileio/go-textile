import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './App'
import Store from './Store'
import * as serviceWorker from './serviceWorker'
import { Provider } from 'mobx-react'
import 'semantic-ui-css/semantic.min.css'

const store = new Store()

const app = (
  <Provider store={store}>
    <App />
  </Provider>
)

ReactDOM.render(app, document.getElementById('root'))

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister()
