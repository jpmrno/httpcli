package httpcli

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc
