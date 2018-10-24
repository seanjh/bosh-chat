async function sleeper (ms = 1000) {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve()
    }, ms)
  })
}

const DEFAULT_WAIT_MS = 1000

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
  for await (const resp of loadMessages()) {
    console.debug(`Received ${resp.status} response.`)
  }
  console.debug('Giving up.')
})()
