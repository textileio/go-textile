const gateway = "https://localhost:9192";
const session = require('electron').remote.session.defaultSession;

let textile = {

  init: function() {
    asticode.loader.init();
    asticode.modaler.init();
    asticode.notifier.init();

    document.addEventListener('astilectron-ready', function() {
      textile.listen()
    })
  },

  pair: function () {
    console.debug("SENDING MESSAGE:", "pair.start");
    /** @namespace astilectron.sendMessage **/
    astilectron.sendMessage({name: "pair.start", payload: ""}, function (message) {
      if (message.name === "error") {
        asticode.notifier.error("Error");
        return
      }

      // populate qr code
      console.debug("GOT QR CODE:", message);
      let qrCode = document.querySelector('.qr-code');
      qrCode.setAttribute('src', "data:image/png;base64," + message.payload.png + "");
      let pairCode = document.querySelector('.confirmation-code');
      pairCode.innerText = message.payload.code
    })
  },

  start: function (pairedID) {
    /** @namespace astilectron.sendMessage **/
    astilectron.sendMessage({name: "sync.start", payload: pairedID}, function (message) {
      if (message.name === "error") {
        asticode.notifier.error("Error")
      }
    })
  },

  listen: function() {
    /** @namespace astilectron.onMessage **/
    astilectron.onMessage(function(message) {
      switch (message.name) {

        case "login.cookie":
          // Setup cookie session for this client
          let expiration = new Date();
          let hour = expiration.getHours();
          hour = hour + 6;
          expiration.setHours(hour);
          /** @namespace ses.cookies **/
          session.cookies.set({
              url: gateway,
              name: message.name,
              value: message.value,
              expirationDate: expiration.getTime(),
              session: true
          }, function (error) {
              // console.log(error);
          });
          break;

        // node and services are ready
        case "sync.ready":
          showGallery(message.html);
          /** @namespace message.pairedID **/
          textile.start(message.pairedID);
          break;

        // new photo from paired peer
        case "sync.data":
          let url = [gateway, "ipfs", message.hash, "thumb"].join("/");
          let $item = $('<div class="grid-item"><img src="' + url + '" /></div>');
          $(".grid").isotope('insert', $item);
          break;

        // start walk-through
        case "onboard.start":
          showOnboarding(1);
          textile.pair();
          break;

        // done on-boarding, we should now have a paired peer
        case "onboard.complete":
          hideOnboarding();
          break
      }
    })
  },
};

function showOnboarding(screen) {
  /** @namespace $ **/
  let ob = $(".onboarding");
  if (ob.hasClass("hidden")) {
    ob.removeClass("hidden")
  }
  $(".onboarding .screen").addClass("hidden");
  $("#ob" + screen).removeClass("hidden")
}

function hideOnboarding() {
  $(".onboarding").addClass("hidden")
}

function showGallery(html) {
  let grid = $(".grid");
  grid.removeClass("hidden");
  grid.html(html);

  // init Isotope
  let $grid = grid.isotope({
    layoutMode: 'cellsByRow',
    itemSelector: '.grid-item',
    cellsByRow: {
      columnWidth: 256,
      rowHeight: 256
    }
  });

  // layout after each image loads
  $grid.imagesLoaded().progress(function() {
    $grid.isotope('layout')
  });

  // reveal items
  let $items = $grid.find('.grid-item');
  $grid.addClass('is-showing-items').isotope('revealItemElements', $items)
}
