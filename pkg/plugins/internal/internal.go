package internal

import (
	"context"
	"errors"
	"net/rpc"
	"strconv"

	hashiplug "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type pluginImplementation struct {
	srv pluginServer
}

// Server implements hashicorp/plugin.Plugin for type checking.
func (*pluginImplementation) Server(*hashiplug.MuxBroker) (interface{}, error) {
	return nil, errors.New("plugins: net/rpc is not used")
}

// Client implements hashicorp/plugin.Plugin for type checking.
func (*pluginImplementation) Client(*hashiplug.MuxBroker, *rpc.Client) (interface{}, error) {
	return nil, errors.New("plugins: net/rpc is not used")
}

// GRPCServer implements hashicorp/plugin.GRPCPlugin for gRPC.
func (p *pluginImplementation) GRPCServer(_ *hashiplug.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServer(s, &p.srv)
	return nil
}

// GRPCClient implements hashicorp/plugin.GRPCPlugin for gRPC.
func (p *pluginImplementation) GRPCClient(_ context.Context, _ *hashiplug.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{
		PluginClient: NewPluginClient(c),
		diskName:     diskName(p.srv.Registration),
	}, nil
}

type pluginServer struct {
	*Registration
}

// Init returns the internal plugin's Registration object.
// (The request object is unused at this time.)
func (p *pluginServer) Init(context.Context, *InitRequest) (*Registration, error) {
	return p.Registration, nil
}

func (*pluginServer) DoHook(_ context.Context, req *HookRequest) (*HookResponse, error) {
	return PluginSideHooks.Do(req.HookName, req.Data)
}

func (*pluginServer) DoFilter(_ context.Context, req *FilterData) (*FilterData, error) {
	return PluginSideFilters.do(req.FilterName, req.Data)
}

// grpcClient is an implementation of hashicorp/plugin.Plugin that talks over gRPC.
type grpcClient struct {
	PluginClient
	diskName string
}

func (m *grpcClient) Init() (*Registration, error) {
	return m.PluginClient.Init(context.Background(), &InitRequest{})
}

func (m *grpcClient) DoHook(hook string, data []byte) ([]byte, error) {
	resp, err := m.PluginClient.DoHook(context.Background(), &HookRequest{
		HookName: hook,
		Data:     data,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, err
}

func (m *grpcClient) DoFilter(filter string, data []byte) ([]byte, error) {
	resp, err := m.PluginClient.DoFilter(context.Background(), &FilterData{
		FilterName: filter,
		Data:       data,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// diskName is the name with which the plugin should be written out to disk.
func diskName(reg *Registration) string {
	conv := func(u uint32) string {
		return strconv.FormatUint(uint64(u), 10)
	}
	return reg.Id + conv(reg.Ver.Major) + "_" + conv(reg.Ver.Minor) + "_" + conv(reg.Ver.Patch)
}
