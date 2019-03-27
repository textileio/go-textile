declare namespace astilectron {
  function onMessage(callback: (message: any) => void)
  function sendMessage(message: any, callback?: (message: any) => void)
}