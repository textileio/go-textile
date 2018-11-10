const session = require('electron').remote.session.defaultSession

let target_dir = ''
let textile = {
  init: function() {
    asticode.loader.init()
    asticode.modaler.init()
    asticode.notifier.init()

    document.addEventListener('astilectron-ready', function() {
      textile.listen()
    })
  },

  listen: function() {
    astilectron.onMessage(function(msg) {
      switch (msg.name) {
        case 'setup':
          showSetup()
          break

        case 'pair':
          showPairing(msg)
          break

        case 'preready':
          showWarmup()
          break

        case 'ready':
          showMain()
          break
      }
    })
  },
}

function focus(view) {
  $('.view').addClass('hidden')
  $('.' + view).removeClass('hidden')
}

function showSetup() {
  focus('setup')
}

function showConfirm(dir) {
  target_dir = dir
  $('.target_dir').text(target_dir)
  focus('confirm')
}

function showPairing(msg) {
  focus('pairing')
  let qrCode = $('.qr-code')
  qrCode.attr('src', 'data:image/png;base64,' + msg.qr)
  qrCode.removeClass('hidden')
  $('.address').text('Address: ' + msg.pk)
}

function showWarmup() {
  focus('warmup')
}

function showMain() {
  focus('main')
}

function rejectTargetFolder() {
  target_dir = ''
  showSetup()
}

function setTargetFolder() {
  astilectron.sendMessage({name: 'openFolderDialog', payload: target_dir}, function (message) {
    if (message.name === 'error') {
      asticode.notifier.error(message)
    }
  })
}

function openFolderDialog() {
  astilectron.showOpenDialog({
    properties: ['openDirectory'],
    title: 'Select a Folder',
  },
  function callback(filePaths) {
    if (filePaths && filePaths.length > 0) {
      showConfirm(filePaths[0])
    }
  });
  return {payload: 'openFolderDialog'};
}


function hideSetup() {
  let setup = $('.setup')
  if (!setup.hasClass('hidden')) {
    setup.addClass('hidden')
  }
}

function showMainO() {
  let main = $('.main')
  if (main.hasClass('hidden')) {
    main.removeClass('hidden')
  }
}

function refresh() {
  astilectron.sendMessage({name: 'refresh'}, function (message) {
    if (message.name === 'error') {
      asticode.notifier.error(message)
    }
  })
  $('.refresh-button').addClass('rotate')
  setTimeout(function () {
    $('.refresh-button').removeClass('rotate')
  }, 500)
}

function renderThreads(threads) {
  threads.forEach(function (thread) {
    addThread(thread)
  })
  if (threads.length > 0) {
    loadFirstThread()
  }
}

function addThread(update) {
  let ul = $('.threads')
  let title = '<h5># ' + update.name + '</h5>'
  $('<li class="thread" id="' + update.id + '" onclick="loadThread(this)">' + title + '</li>').appendTo(ul)
  if (ul.children().length === 1) {
    loadFirstThread()
  }
}

function removeThread(update) {
  let ul = $('.threads')
  let active = $('.thread.active').attr('id') === update.id
  $('#' + id).remove()
  if (active) {
    if (ul.children().length > 0) {
      loadFirstThread()
    } else {

    }
  }
}

function loadFirstThread() {
  setTimeout(function () {
    $('.threads li').first().click()
  }, 0)
}

function loadThread(el) {
  let $el = $(el)
  let id = $el.attr('id')
  $('.thread.active').removeClass('active')
  $el.addClass('active')
  astilectron.sendMessage({name: 'thread.load', payload: id}, function (message) {
    if (message.name === 'error') {
      asticode.notifier.error(message)
      return
    }
    showGrid(id, message.payload.html)
  })
}

function showGrid(threadId, html) {
  clearGrid()
  $('.message').addClass('hidden')
  let grid = $('<div class="grid" data-thread-id="' + threadId + '"></div>')
  grid.appendTo($('.content'))

  grid.html(html)
  let $grid = grid.isotope({
    layoutMode: 'cellsByRow',
    itemSelector: '.grid-item',
    cellsByRow: {
      columnWidth: 192,
      rowHeight: 192
    },
    transitionDuration: '0.2s',
    hiddenStyle: {
      opacity: 0,
      transform: 'scale(0.9)'
    },
    visibleStyle: {
      opacity: 1,
      transform: 'scale(1)'
    }
  })

  // layout after each image loads
  $grid.imagesLoaded().progress(function() {
    if ($grid.data('isotope')) {
      $grid.isotope('layout')
    }
  })

  // reveal items
  let $items = $grid.find('.grid-item')
  $grid.addClass('is-showing-items').isotope('revealItemElements', $items)
}

function clearGrid() {
  $('.grid').remove()
}

function addPhoto(update) {
  let grid = $('.grid')
  if (!grid || update.thread_id !== grid.data('thread-id')) {
    return
  }
  let photo = fileURL(update, 'photo')
  let small = fileURL(update, 'small')
  let meta = fileURL(update, 'meta')
  let img = '<img src="' + small + '" />'
  let $item = $('<div id="' + update.block.id + '" class="grid-item" '
    + 'ondragstart="imageDragStart(event);" draggable="true" '
    + 'data-url="' + photo + '" data-meta="' + meta + '">' + img + '</div>')
  grid.isotope().prepend($item).isotope('prepended', $item)
}

function ignore(update) {
  let grid = $('.grid')
  if (!grid || update.thread_id !== grid.data('thread-id')) {
    return
  }
  if (!update.block.data_id) {
    return
  }
  let parts = update.block.data_id.split('-')
  if (parts.length !== 2) {
    return
  }
  grid.isotope('remove', $('#' + parts[1])).isotope('layout')
}

function fileURL(update, path) {
  return [textile.gateway, 'ipfs', update.block.data_id, path].join('/') + '?block=' + update.block.id
}

function login(data) {
  textile.gateway = data.gateway
  let expiration = new Date()
  let hour = expiration.getHours()
  hour = hour + 6
  expiration.setHours(hour)
  session.cookies.set({
    url: data.gateway,
    name: data.name,
    value: data.value,
    expirationDate: expiration.getTime(),
    session: true
  }, function (err) {
    // console.error(err)
  })
}