package ipfs

import "gx/ipfs/QmXporsyf5xMvffd2eiTDoq85dNpYUynGJhfabzDjwP8uR/go-ipfs/commands"

func ConnectedPeers(ctx commands.Context) ([]string, error) {
	args := []string{"swarm", "peers"}
	req, cmd, err := NewRequest(ctx, args)
	if err != nil {
		return nil, err
	}
	res := commands.NewResponse(req)
	cmd.Run(req, res)
	if res.Error() != nil {
		return nil, res.Error()
	}
	return *res.Output().(*[]string), nil
}
