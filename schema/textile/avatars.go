package textile

var Avatars = `
{
  "pin": true,
  "links": {
    "small": {
      "use": ":file",
      "pin": true,
      "plaintext": true,
      "mill": "/image/resize",
      "opts": {
        "width": "320",
        "quality": "75"
      }
    },
    "thumb": {
      "use": ":file",
      "pin": true,
      "plaintext": true,
      "mill": "/image/resize",
      "opts": {
        "width": "100",
        "quality": "75"
      }
    }
  }
}
`
