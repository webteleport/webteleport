//go:build js

// Command wasm is a js/wasm binary that exposes the webteleport Listen
// function to JavaScript via syscall/js.
//
// Build with:
//
//	GOOS=js GOARCH=wasm go build -o main.wasm ./cmd/wasm
//
// Load in a browser page alongside wasm_exec.js (shipped with the Go toolchain):
//
//	<script src="wasm_exec.js"></script>
//	<script>
//	  const go = new Go();
//	  WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
//	    .then(result => go.run(result.instance));
//	</script>
//
// After the module is running, call:
//
//	const ln = await webteleportListen("wss://relay.example.com");
//	console.log("assigned address:", ln.addr);
//	while (true) {
//	  const conn = await ln.accept();
//	  handleConn(conn);
//	}
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"syscall/js"

	"github.com/webteleport/webteleport"
)

func main() {
	js.Global().Set("webteleportListen", js.FuncOf(jsListen))
	// Block forever so the wasm module stays alive.
	select {}
}

// jsListen is the JS-callable binding for webteleport.Listen.
//
// Signature: webteleportListen(relayAddr: string) -> Promise<Listener>
//
// Listener fields:
//
//	addr    string   – host[:port] assigned by the relay server
//	accept() -> Promise<Conn>
//	close()
//
// Conn fields:
//
//	localAddr   string
//	remoteAddr  string
//	read(n?: number) -> Promise<Uint8Array | null>  (null signals EOF)
//	write(data: Uint8Array) -> Promise<number>
//	close()
func jsListen(this js.Value, args []js.Value) any {
	return newPromise(func() (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), fmt.Errorf("webteleportListen: relayAddr argument required")
		}
		relayAddr := args[0].String()
		ln, err := webteleport.Listen(context.Background(), relayAddr)
		if err != nil {
			return js.Undefined(), err
		}
		return listenerToJS(ln), nil
	})
}

// listenerToJS wraps a net.Listener as a plain JS object.
func listenerToJS(ln net.Listener) js.Value {
	var funcs []js.Func

	acceptFn := js.FuncOf(func(this js.Value, _ []js.Value) any {
		return newPromise(func() (js.Value, error) {
			conn, err := ln.Accept()
			if err != nil {
				return js.Undefined(), err
			}
			return connToJS(conn), nil
		})
	})
	funcs = append(funcs, acceptFn)

	closeFn := js.FuncOf(func(this js.Value, _ []js.Value) any {
		err := ln.Close()
		for _, f := range funcs {
			f.Release()
		}
		if err != nil {
			return js.ValueOf(err.Error())
		}
		return js.Null()
	})
	funcs = append(funcs, closeFn)

	obj := js.Global().Get("Object").New()
	obj.Set("addr", ln.Addr().String())
	obj.Set("accept", acceptFn)
	obj.Set("close", closeFn)
	return obj
}

// connToJS wraps a net.Conn as a plain JS object.
func connToJS(conn net.Conn) js.Value {
	var funcs []js.Func

	readFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		n := 4096
		if len(args) > 0 && args[0].Type() == js.TypeNumber {
			n = args[0].Int()
		}
		return newPromise(func() (js.Value, error) {
			buf := make([]byte, n)
			nr, err := conn.Read(buf)
			if err == io.EOF {
				return js.Null(), nil
			}
			if err != nil {
				return js.Undefined(), err
			}
			dst := js.Global().Get("Uint8Array").New(nr)
			js.CopyBytesToJS(dst, buf[:nr])
			return dst, nil
		})
	})
	funcs = append(funcs, readFn)

	writeFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		return newPromise(func() (js.Value, error) {
			if len(args) < 1 {
				return js.Undefined(), fmt.Errorf("write: Uint8Array argument required")
			}
			src := args[0]
			data := make([]byte, src.Length())
			js.CopyBytesToGo(data, src)
			nw, err := conn.Write(data)
			if err != nil {
				return js.Undefined(), err
			}
			return js.ValueOf(nw), nil
		})
	})
	funcs = append(funcs, writeFn)

	closeFn := js.FuncOf(func(this js.Value, _ []js.Value) any {
		err := conn.Close()
		for _, f := range funcs {
			f.Release()
		}
		if err != nil {
			return js.ValueOf(err.Error())
		}
		return js.Null()
	})
	funcs = append(funcs, closeFn)

	obj := js.Global().Get("Object").New()
	obj.Set("localAddr", conn.LocalAddr().String())
	obj.Set("remoteAddr", conn.RemoteAddr().String())
	obj.Set("read", readFn)
	obj.Set("write", writeFn)
	obj.Set("close", closeFn)
	return obj
}

// newPromise returns a JS Promise that executes fn in a new goroutine.
// fn should return (resolve value, error); on error the promise rejects
// with a JS Error whose message is err.Error().
func newPromise(fn func() (js.Value, error)) js.Value {
	var handler js.Func
	handler = js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve, reject := args[0], args[1]
		go func() {
			defer handler.Release()
			val, err := fn()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New(err.Error()))
				return
			}
			resolve.Invoke(val)
		}()
		return nil
	})
	return js.Global().Get("Promise").New(handler)
}
