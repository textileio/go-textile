//  Created by react-native-create-bridge

import { NativeModules } from 'react-native'

const { TextileIPFS } = NativeModules

export default {
  exampleMethod () {
    return TextileIPFS.exampleMethod()
  },

  EXAMPLE_CONSTANT: TextileIPFS.EXAMPLE_CONSTANT
}
