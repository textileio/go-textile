let textile = {
  init: function() {
    // Init
    asticode.loader.init();
    asticode.modaler.init();
    asticode.notifier.init();

    this.gateway = "http://localhost"

    // Wait for astilectron to be ready
    document.addEventListener('astilectron-ready', function() {
      // Listen
      textile.listen();

      // let pairNewDevice = function () {
      //   // Get a thumbnail (i know it exists)
      //   astilectron.sendMessage({name: "peer.qr", payload: ""}, function (message) {
      //     // Check error
      //     console.log(message)
      //     if (message.name === "error") {
      //       asticode.notifier.error("Error");
      //       return
      //     }
      //     let qrCode = document.querySelector('.peerQRCode');
      //     qrCode.setAttribute('src', "data:image/png;base64," + message.payload.png + "")
      //     let pairCode = document.querySelector('.pairCode');
      //     pairCode.innerText = message.payload.code
      //     let modal = document.querySelector('.modal');
      //     modal.classList.toggle('modal-open');
      //   })
      // }
      // setupQRModal(pairNewDevice)
    })
  },
  pair: function () {
      astilectron.sendMessage({name: "pairing.start", payload: ""}, function (message) {
        // Check error
        if (message.name === "error") {
          asticode.notifier.error("Error");
          return
        }
        let qrCode = document.querySelector('.qr-code');
        qrCode.setAttribute('src', "data:image/png;base64," + message.payload.png + "")
        let pairCode = document.querySelector('.confirmation-code');
        pairCode.innerText = message.payload.code
      })
  },
  listen: function() {
    astilectron.onMessage(function(message) {
      switch (message.name) {
        case "ready":
          /*
          When textile-go has finished all the startup processes, this fires
           */
          this.gateway = message.gateway;
          initPhotoSwipeFromDOM('.gallery');
        case "sync.new":
          /*
          When textile-go receives a new photo from the mobile device, this should fire...
           */
          console.log(message)
          let image = $('<img src="https://gateway.ipfs.io/ipfs/QmR8mGCutBWDPBc9zdpPZPoRYRAJS7BMZhJtqHeFtJp2ma/thumb.jpg"  itemprop="thumbnail" alt="Image description"/>');
          image.attr('src', this.gateway + "/ipfs/" + message.hash + "/thumb")
            console.log(this.gateway + "/ipfs/" + message.hash + "/thumb")
          let link = $('<a href="https://farm2.staticflickr.com/1043/5186867718_06b2e9e551_b.jpg" itemprop="contentUrl" data-size="964x1024"></a>')
          link.attr('href', this.gateway + "/ipfs/" + message.hash + "/photo")
          let caption = $('<figcaption itemprop="caption description">Image caption  1</figcaption>')
          caption.text(message.timestamp)
          let figure = $('<figure itemprop="associatedMedia" itemscope itemtype="http://schema.org/ImageObject"></figure>')
          link.append(image)
          figure.append(link)
          figure.append(caption)

          $('.gallery').append(figure);
          // need to rewrite this, shouldn't do it this way, just trying to get to a working state
          initPhotoSwipeFromDOM('.gallery');
          break;
        case "onboard.complete":
          /*
          User is in onboarding screens, stuck on the QR code screen.
          This tells the UI that pairing is complete, so move past that screen
           */
          console.log("pairing complete", message);
          nextScreen();
          break;
        case "onboard":
          /*
          This kicks of the UI for explaining the onboarding process and presenting the QR code
           */
          console.log("onboarding")
          $('.onboarding').removeClass('hidden');
          $('.onboarding').find('.screen').first().removeClass('hidden');
          $('.onboarding').find('.screen').first().find('.next').click(nextScreen);
          textile.pair();
        case "new.log":
          console.log("log:", message)
          break;
      }
    });
  },
};

var nextScreen = function () {
  $('.onboarding').find('.screen').first().remove();
  $('.onboarding').find('.screen').first().removeClass('hidden');
  setTimeout(function(){ // avoids multiple clicks on button
    $('.onboarding').find('.screen').first().find('.next').click(nextScreen)
  },300);
}
