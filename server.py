#!/usr/bin/env python
import SimpleHTTPServer
import SocketServer

PORT = 8000

class MyHTTPRequestHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_my_headers()

        SimpleHTTPServer.SimpleHTTPRequestHandler.end_headers(self)

    def send_my_headers(self):
        self.send_header("Access-Control-Allow-Origin", "*")


if __name__ == '__main__':
    httpd = SocketServer.TCPServer(("", PORT), MyHTTPRequestHandler)
    print "serving at port", PORT
    httpd.serve_forever()