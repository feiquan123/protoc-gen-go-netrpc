package gengonetrpcprev

const (
	tmpService = `
{{$root := .}}

// encode format
type Encode string

const (
	EncodeJson = Encode("json")
	EncodeGob = Encode("gob")
)

// interface defined
type {{.ServiceName}}Interface interface {
	{{- range $_, $m := .MethodList}}
	{{$m.MethodName}} (*{{$m.InputTypeName}}, *{{$m.OutputTypeName}}) error
	{{- end}}
}

var registerOnce sync.Once
// register rpc service
func Register{{.ServiceName}}(srv *rpc.Server, x {{.ServiceName}}Interface) error{
	var err error
	registerOnce.Do(func() {
		err = srv.RegisterName("{{.ServiceName}}",x)
	})
	if err != nil{
		return err
	}
	return nil
}

// run service func
func GenRun{{.ServiceName}}(srv *rpc.Server, x {{.ServiceName}}Interface,network, address string, encode Encode) (func ()error, error){
	err := Register{{.ServiceName}}(srv, x)
	if err != nil{
		return nil, fmt.Errorf("register error: %s", err)
	}

	listener, err := net.Listen(network, address)
	if err !=nil{
		return nil, fmt.Errorf("listen %s error: %s", network, address)
	}

	return func() error {
		for {
			conn, err := listener.Accept()
			if err != nil{
				return err
			}
			go func(encode Encode){
				switch encode{
				case EncodeJson:
					rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
				case EncodeGob:
					rpc.ServeConn(conn)
				default:
					rpc.ServeConn(conn)
				}
			}(encode)
		}
	}, nil
}


// client struct 
type {{.ServiceName}}Client struct{
	*rpc.Client
}

var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client) (nil)

func Dial{{.ServiceName}}(network, address string, encode Encode) (*{{.ServiceName}}Client, error){
	conn, err := net.Dial(network, address)
	if err != nil{
		return nil, err
	}

	var c *rpc.Client
	switch encode{
	case EncodeJson:
		c = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	case EncodeGob:
		c = rpc.NewClient(conn)
	default:
		c = rpc.NewClient(conn)
	}

	return &{{.ServiceName}}Client{
		Client: c,
	}, nil
}

{{range $_, $m := .MethodList}}
func (p *{{$root.ServiceName}}Client) {{$m.MethodName}} (in *{{$m.InputTypeName}}, out *{{$m.OutputTypeName}}) error {
	return p.Client.Call("{{$root.ServiceName}}.{{$m.MethodName}}", in, out)
}
{{end}}
`
)

type ServiceSpec struct {
	ServiceName string
	MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
	MethodName     string
	InputTypeName  string
	OutputTypeName string
}
