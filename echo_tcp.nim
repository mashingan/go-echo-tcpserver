import net, asyncnet, asyncdispatch

const port = Port 3000

proc echoBack(client: AsyncSocket) {.async.} =
  echo "accepted connection."
  let line = await client.recvLine()
  echo "getting ", line
  echo "sending it back"
  await client.send("echo> " & line & "\c\L")
  client.close

proc serve() {.async.} =
  var server = newAsyncSocket()
  server.setSockOpt(OptReuseAddr, true)
  server.bindAddr port
  server.listen
  echo "server listening on 127.0.0.1 with port ", port

  while true:
    let client = await server.accept()
    asyncCheck echoBack(client)


asyncCheck serve()
runForever()
