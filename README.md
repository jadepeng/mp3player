# mp3player
a golang websocket mp3 player

the player start a websocket server `ws://localhost:8080/ws`, you can use `--addr=localhost:8081` to custom ws server port

javascript client can send `{"command":"play","arg":"mp3_file_path"}` to let server play mp3 file
