const DEFAULT_WAIT_MS = 1000
const MAX_MESSAGES = 1000
const KEY_ENTER = 'Enter'

function postData (url = '', data = {}) {
  return fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json; charset=utf-8'
    },
    body: JSON.stringify(data)
  })
}

function getMessage () {
  return document.getElementById('message')
}

function handleMessageKey (e) {
  if (e.key === KEY_ENTER) {
    e.preventDefault()
    submitMessage()
  }
}

function setupSubmit () {
  const msgElem = getMessage()
  msgElem.onkeyup = handleMessageKey
  const sendElem = document.getElementById('send')
  sendElem.onclick = submitMessage
}

async function submitMessage () {
  const elem = getMessage()
  const message = elem.value
  console.debug(`Submitting message "${message}"`)
  elem.value = ''
  await postData('/messages/', { body: message })
}

function appendMessages (messages = [], buff = []) {
  const content = messages.map(m => m.body)
  if (content.length > 0) {
    console.debug(`Appending ${content.length} total messages.`)
    const total = content.length + buff.length
    const deficit = MAX_MESSAGES - total
    if (deficit < 0) {
      buff = buff.slice(0 - deficit).concat(content)
    } else {
      buff = buff.concat(content)
    }
  } else {
    console.debug('No messages')
  }
  return buff
}

function getChat () {
  return document.getElementById('chat')
}

function writeMessages (buff) {
  const elem = getChat()
  elem.value = buff.join('\n')
}

async function sleeper (ms = 1000) {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve()
    }, ms)
  })
}

async function * loadMessages () {
  let failure = 0
  let ms = DEFAULT_WAIT_MS
  while (true && failure < 10) {
    const wait = sleeper(ms)
    const resp = await fetch('/messages/')
    yield resp
    if (resp.status !== 200) {
      failure++
      ms *= 2
    } else {
      failure = 0
      ms = DEFAULT_WAIT_MS
    }
    await wait
  }
}

(async function main () {
  setupSubmit()

  let buff = []
  for await (const resp of loadMessages()) {
    console.debug(`Received ${resp.status} response.`)
    buff = appendMessages(await resp.json(), buff)
    writeMessages(buff)
  }
  console.debug('Giving up.')
})()
